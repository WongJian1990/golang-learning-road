package main

import (
	"logagent/config"
	"logagent/log"

	"github.com/astaxie/beego/logs"
)

func main() {
	err := config.LogAgentCtx.InitLogAgentContext("ini", "config/app.conf")
	if err != nil {
		logs.Error(err)
	}
	log.InitLogs(config.LogAgentCtx.LogsConfig)
}
