package server

import (
	"errors"
	"logmanager/config"

	"github.com/astaxie/beego/logs"
	elastic "gopkg.in/olivere/elastic.v2"
)

type EsClient struct {
	client *elastic.Client
}

func NewEsClient(config *config.EsConfig) (*EsClient, error) {
	if config == nil {
		return nil, errors.New("NewEsClient: invalid config")
	}
	logs.Debug("NewEsClient URL: ", config.Url)
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(config.Url))
	return &EsClient{client}, err
}

func (p *EsClient) Run() {
	for {
		// logs.Debug("EsClient Run .")
		select {
		case msg := <-kafkaToEsChan:
			// logs.Debug("Kafka")
			_, err := p.client.Index().
				Index(msg.Topic).
				Type(msg.Topic).
				BodyJson(msg).Do()
			if err != nil {
				logs.Error("EsClient::Run: ", err)
			}
			// logs.Debug("EsClient::Add Messege success")
		}
	}
}
