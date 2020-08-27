package monitor

import (
	"errors"
	"fmt"
	"logagent/config"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type KafkaProducer struct {
	client sarama.AsyncProducer
}

func NewKafkaProducer(config *config.KafkaConfig) (*KafkaProducer, error) {
	if config == nil {
		return nil, errors.New("NewKafkaProducer:KafkaConfig invalid")
	}

	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	client, err := sarama.NewAsyncProducer(config.EndPoints, conf)
	if err != nil {
		return nil, fmt.Errorf("NewKafkaProducer: %s", err)
	}

	return &KafkaProducer{client}, nil
}

//Run 启动Kafka生产监测
func (p *KafkaProducer) Run() {
	defer p.client.AsyncClose()
	for {
		//异常模式处理Errors/Successes 否则可能死锁，导致消费者接收不到数据
		errors := p.client.Errors()
		success := p.client.Successes()
		select {
		case msg := <-tailfToKafkaChan:
			kmsg := &sarama.ProducerMessage{
				Topic: msg.Topic,
				Value: sarama.StringEncoder(msg.Value),
			}
			p.client.Input() <- kmsg

			// _, _, err := p.client.SendMessage(kmsg)
			// if err != nil {
			// 	logs.Error("KafkaProducer::Run: SendMessage ", err)
			// }

		case <-success:
			continue
		case <-errors:
			logs.Error("KafkaProducer::Run: SendMessage error")
		}
	}
}
