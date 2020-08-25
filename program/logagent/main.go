package main

import (
	"logagent/config"
	"logagent/log"
	"logagent/monitor"

	"github.com/astaxie/beego/logs"
)

func main() {
	err := config.LogAgentCtx.InitLogAgentContext("ini", "config/app.conf")
	if err != nil {
		logs.Error(err)
	}
	//日志初始化
	log.InitLogs(config.LogAgentCtx.LogsConfig)

	client, err := monitor.NewEtcdClient(config.LogAgentCtx.EtcdConfig)
	if err != nil {
		logs.Error(err)
	}
	//启动Etcd监测服务
	client.Start()
	defer client.ShutDown()

	//启动Tail监测服务
	tailManager := monitor.NewTailfManager()
	tailManager.Start()
	defer tailManager.ShutDown()

	//监测运行
	monitor.Run()
}
