package log

import (
	"encoding/json"
	"logagent/config"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func InitLogs(logConfig *config.LogsConfig) error {
	LogsFileJson, err := json.Marshal(*logConfig)
	if err != nil {
		logs.Error("InitLogs: ", err)
		return err
	}

	err = beego.SetLogger(logs.AdapterFile, string(LogsFileJson))
	beego.SetLogFuncCall(true)
	if err != nil {
		logs.Error("InitLogs: ", err)
		return err
	}
	logs.Info("InitLogs Sucess")
	return nil
}
