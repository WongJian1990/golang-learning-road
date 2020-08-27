package main

import (
	"logmanager/config"
	"logmanager/log"
	"logmanager/server"

	"github.com/astaxie/beego/logs"
)

func main() {
	err := config.LogAgentCtx.InitLogAgentContext("ini", "config/app.conf")
	if err != nil {
		logs.Error(err)
		return
	}
	//日志初始化
	log.InitLogs(config.LogAgentCtx.LogsConfig)

	client, err := server.NewEtcdClient(config.LogAgentCtx.EtcdConfig)
	if err != nil {
		logs.Error(err)
		return
	}
	//启动Etcd监测服务
	client.Start()
	defer client.ShutDown()

	//添加kafka管理pod
	m, err := server.NewKafkaConsumerManager(config.LogAgentCtx.KafkaConfig)
	if err != nil {
		logs.Error(err)
		return
	}
	server.AddPod(m)

	//添加es客户端pod
	e, err := server.NewEsClient(config.LogAgentCtx.EsConfig)
	if err != nil {
		logs.Error(err)
		return
	}
	server.AddPod(e)

	server.Run()
}
