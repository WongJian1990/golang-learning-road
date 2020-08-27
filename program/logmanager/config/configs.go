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
	EsConfig    *EsConfig
	KafkaConfig *KafkaConfig
}

type LogsConfig struct {
	LogsFileName string `json:"filename"`
	LogsLevel    int    `json:"level"`
}

type EtcdConfig struct {
	EndPoints          []string
	DialTimeOut        int
	WatchKeyPrefix     string
	StatusKey          string
	Etcd2KakfaChanSize int
}

type EsConfig struct {
	Url string
}

type KafkaConfig struct {
	EndPoints        []string
	Kafka2EsChanSize int
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
	//加载配置文件
	err := beego.LoadAppConfig(adapter, fileName)
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	//解析Log配置
	LogsConfig, err := initLogsConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	ctx.LogsConfig = LogsConfig

	//解析Etcd配置
	EtcdConfig, err := initEtcdConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	ctx.EtcdConfig = EtcdConfig

	//解析Kafka配置
	KafkaConfig, err := initKafkaConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	ctx.KafkaConfig = KafkaConfig

	//解析Es配置
	EsConfig, err := initEsConfig()
	if err != nil {
		logs.Error("InitLogAgentContext: ", err)
		return err
	}
	ctx.EsConfig = EsConfig

	return nil
}

//转为endpoints
func toEndPoints(serial string) ([]string, error) {
	serial = strings.TrimSpace(serial)
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

//解析Ectd相关配置
func initEtcdConfig() (*EtcdConfig, error) {
	EndPointsSerial := beego.AppConfig.String("etcd_end_points")
	if len(EndPointsSerial) == 0 {
		return nil, fmt.Errorf("initEtcdConfig: EndPoints Array is nil")
	}
	endPoints, err := toEndPoints(EndPointsSerial)
	if err != nil {
		return nil, fmt.Errorf("initEtcdConfig: %s", err)
	}
	timeout, err := beego.AppConfig.Int("etcd_dial_timeout")
	if err != nil {
		logs.Warn("Etcd Timeout is error, 5 as default")
		timeout = 5
	}
	watchKeyPrefix := beego.AppConfig.String("etcd_watch_key_prefix")
	if len(watchKeyPrefix) == 0 {
		logs.Warn("Etcd WatchKeyPrefix is empty, none as default")
		watchKeyPrefix = "none"
	}

	statusKey := beego.AppConfig.String("etcd_watch_status_key")
	if len(statusKey) == 0 {
		logs.Warn("Etcd WatchKeyPrefix is empty, private/status as default")
		statusKey = "private/status"
	}

	etcd2KafkaChanSize, err := beego.AppConfig.Int("etcd_to_kafka_chan_size")
	if err != nil {
		logs.Warn("Etcd Etcd2TailChanSize is error,20 as default")
		etcd2KafkaChanSize = 20
	}

	return &EtcdConfig{
		endPoints,
		timeout,
		watchKeyPrefix,
		statusKey,
		etcd2KafkaChanSize,
	}, nil
}

//解析Kafka相关配置
func initKafkaConfig() (*KafkaConfig, error) {
	EndPointsSerial := beego.AppConfig.String("kafka_end_points")
	if len(EndPointsSerial) == 0 {
		return nil, fmt.Errorf("initKafkaConfig: EndPoints Array is nil")
	}
	endPoints, err := toEndPoints(EndPointsSerial)
	if err != nil {
		return nil, fmt.Errorf("initKafkaConfig: %s", err)
	}

	kafka2EsChanSize, err := beego.AppConfig.Int("etcd_to_kafka_chan_size")
	if err != nil {
		logs.Warn("Etcd Kafka2EsChanSize is error,20 as default")
		kafka2EsChanSize = 20
	}
	return &KafkaConfig{EndPoints: endPoints, Kafka2EsChanSize: kafka2EsChanSize}, nil
}

//解析Es相关配置
func initEsConfig() (*EsConfig, error) {
	url := beego.AppConfig.String("es_server_url")
	if len(url) == 0 {
		logs.Warn("initTailfConfig es server url is empty, localhost:9200 as default")
		url = "localhost:9200"
	}

	return &EsConfig{url}, nil
}

//解析日志相关配置
func initLogsConfig() (*LogsConfig, error) {
	LogsFileName := beego.AppConfig.String("logs_filename")
	if len(LogsFileName) == 0 {
		logs.Warn("LogsFileName is empty, logagent.log as default")
		LogsFileName = "logagent.log"
	}
	LogsLevel := beego.AppConfig.String("logs_level")
	Level, err := toLogsLevel(LogsLevel)
	if err != nil {
		logs.Error(err)
		logs.Warn("LogsLevel is invalid, debug as default")
		Level = beego.LevelDebug
	}

	return &LogsConfig{LogsFileName, Level}, nil
}

//转为LogsLevel
func toLogsLevel(LogsLevel string) (level int, err error) {
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
