package main

import (
	"encoding/json"
	"logagent/config"
	"os"
	"time"

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

func Run() {
	index := 1
	for {
		logs.Debug("For Test, ", index)
		index++
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) != 2 {
		logs.Error("Arg: Invalid")
	}
	filename := os.Args[1]
	config := &config.LogsConfig{LogsFileName: filename, LogsLevel: beego.LevelDebug}
	err := InitLogs(config)
	if err != nil {
		return
	}
	Run()
}
