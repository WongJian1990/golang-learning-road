package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"logmanager/config"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type KafkaConsumerManager struct {
	config  *config.KafkaConfig
	clients map[string]*KafkaConsumer
	ctx     context.Context
	cancel  context.CancelFunc
}
type KafkaConsumer struct {
	client sarama.Consumer
	topic  string
	ctx    context.Context
	cancel context.CancelFunc
	// exit   bool
}

func NewKafkaConsumerManager(config *config.KafkaConfig) (*KafkaConsumerManager, error) {
	if config == nil {
		return nil, errors.New("NewKafkaConsumerManager:KafkaConfig invalid")
	}

	ctx, cancel := context.WithCancel(context.Background())
	kafkaToEsChan = make(chan *Msg, config.Kafka2EsChanSize)
	return &KafkaConsumerManager{config, make(map[string]*KafkaConsumer, 0), ctx, cancel}, nil

}

func (m *KafkaConsumerManager) Run() {
	for {
		select {
		case pod := <-etcdPodChan:
			var msg Msg
			err := json.Unmarshal([]byte(pod.value), &msg)
			if err != nil {
				logs.Error("KafkaConsumerManager::Run: Unmarshal ", err)
				continue
			}
			switch pod.kvt {
			case "PUT":
				if m.clients[msg.Topic] != nil {
					logs.Info("KafkaConsumerManager::Run: consumer[", msg.Topic, "] aready run")
					continue
				}
				con, err := m.NewKafkaConsumer(msg.Topic)
				if err != nil {
					logs.Error("KafkaConsumerManager::Run: ", err)
					continue
				}
				m.clients[msg.Topic] = con
				go con.start()

			case "DELETE":
				if m.clients[msg.Topic] != nil {
					m.clients[msg.Topic].shutdown()
					delete(m.clients, msg.Topic)
				}
			}
		case <-m.ctx.Done():
			break
		}
	}
}

func (m *KafkaConsumerManager) NewKafkaConsumer(topic string) (*KafkaConsumer, error) {

	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true

	client, err := sarama.NewConsumer(m.config.EndPoints, conf)
	if err != nil {
		return nil, fmt.Errorf("NewKafkaConsumer: %s", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{client, topic, ctx, cancel}, nil
}

//Run 启动Kafka消费
func (p *KafkaConsumer) start() {
	defer p.client.Close()
	//异常模式处理Errors 否则可能死锁，导致消费者接收不到数据
	partitionList, err := p.client.Partitions(p.topic)
	if err != nil {
		logs.Error("KafkaConsumer::Partitions: err", err)
		return
	}
	//logs.Debug("PartitionList Length: ", len(partitionList))
	for partition := range partitionList {
		pc, err := p.client.ConsumePartition(p.topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			//logs.Error("Consumer::ConsumePartition: ", err)
			continue
		}
		defer pc.AsyncClose()
		for {
			errors := pc.Errors()
			select {
			case msg := <-pc.Messages():
				value := fmt.Sprintf("[msg offset: %d, partition: %d, timestamp: %s] message: %s\n",
					msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
				kmsg := &Msg{p.topic, value}
				kafkaToEsChan <- kmsg
			case <-errors:
				logs.Error("KafkaConsumer::ConsumerPartition ", errors)
			case <-p.ctx.Done():
				return
			}
		}
	}
}

func (p *KafkaConsumer) shutdown() {
	p.cancel()
	// p.exit = true
}
