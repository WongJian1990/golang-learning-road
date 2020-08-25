package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LogAgentContext struct {
	LogsConfig  *LogsConfig
	EtcdConfig  *EtcdConfig
	TailfConfig *TailfConfig
	KafkaConfig *KafkaConfig
}

type LogsConfig struct {
	LogsFileName string `json:"filename"`
	LogsLevel    int    `json:"level"`
}

type EtcdConfig struct {
	EndPoints   []string
	DialTimeOut int
}

type TailfConfig struct {
}

type KafkaConfig struct {
}

//LogAgentCtx export
var LogAgentCtx *LogAgentContext

//初始化
func init() {
	LogAgentCtx = NewLogAgentContext()
}

//NewLogAgentContext 构造日志收集配置上下文
func NewLogAgentContext() *LogAgentContext {
	return &LogAgentContext{}
}

//InitLogAgentContext 获取配置上下文信息
func (ctx *LogAgentContext) InitLogAgentContext(adapter, fileName string) error {
	err := beego.LoadAppConfig(adapter, fileName)
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	LogsConfig, err := initLogsConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	ctx.LogsConfig = LogsConfig
	EtcdConfig, err: = initEtcdConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	ctx.EtcdConfig=EtcdConfig

	err = initTailfConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}

	err = initKafkaConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}

	return nil
}

func toEndPoints(serial string) ([]string, error) {
	endpoints := strings.Split(serial, ",")
	if len(endpoints) == 0 {
		return nil, errors.New("Invalid EndPoints")
	}
	for _, ep := range endpoints {
		addrandport := strings.Split(ep, ":")
		if len(addrandport) != 2 {
			return nil, errors.New("Error EndPoint Format")
		}
		_, err := strconv.Atoi(addrandport[1])
		if err != nil {
			return nil, errors.New("Error EndPoint Port Error")
		}
	}
	return endpoints, nil
}
func initEtcdConfig() (*EtcdConfig, error) {
	EndPointsSerial := beego.AppConfig.String("etcd_end_points")
	if len(EndPointsSerial) == 0 {
		return nil, fmt.Errorf("initEtcdConfig: EndPoints Array is nil")
	}
	endPoints, err := toEndPoints(EndPointsSerial)
	if err != nil {
		return nil, err
	}
	timeout, err := beego.AppConfig.Int("etcd_dial_timeout")
	if err != nil {
		return nil, errors.New("intEtcdConfig: DialTimout error")
	}
	return &EtcdConfig{
		endPoints,
		timeout,
	}, nil
}

func initKafkaConfig() error {
	return nil
}

func initTailfConfig() error {
	return nil
}

func initLogsConfig() (*LogsConfig, error) {
	LogsFileName := beego.AppConfig.String("logs_filename")
	if len(LogsFileName) == 0 {
		logs.Warn("LogsFileName is empty, logagent.log as default")
		LogsFileName = "logagent.log"
	}
	LogsLevel := beego.AppConfig.String("logs_level")
	Level, err := checkLogsLevelValid(LogsLevel)
	if err != nil {
		logs.Error(err)
		logs.Warn("LogsLevel is invalid, debug as default")
		Level = beego.LevelDebug
	}

	return &LogsConfig{LogsFileName, Level}, nil
}

func checkLogsLevelValid(LogsLevel string) (level int, err error) {
	switch strings.ToLower(LogsLevel) {
	case "error":
		level = beego.LevelError
	case "warn":
		level = beego.LevelWarning
	case "info":
		level = beego.LevelInformational
	case "debug":
		level = beego.LevelDebug
	default:
		err = fmt.Errorf("LogsLevel: not support level:%s", err)
	}
	return
}
