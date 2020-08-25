package monitor

type EtcdToTail struct {
	//键值
	key,
	//值
	value,
	//动作类型[PUT|DELETE]
	kvt string
}

var etcdToTailChan chan *EtcdToTail
