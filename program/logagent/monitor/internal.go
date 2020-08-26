package monitor

type EtcdToTail struct {
	//键值
	key,
	//值
	value,
	//动作类型[PUT|DELETE]
	kvt string
}

type Msg struct {
	//主题
	Topic string `json:"topic"`
	//消息
	Value string `json:"message"`
}

var etcdToTailChan chan *EtcdToTail

var tailfToKafkaChan chan *Msg
