package server

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"logmanager/config"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	etcd_client "go.etcd.io/etcd/clientv3"
)

const (
	//EtcdInitial Etcd客户端初始状态
	EtcdInitial = 0
	//EtcdStart 启动Etcd客户端监测服务
	EtcdStart = 1
	//EtcdStop 暂停Etcd客户端监测服务
	EtcdStop = 2
	//EtcdShutDown 关闭Etcd客户端监测服务
	EtcdShutDown = 3
)

var wg *sync.WaitGroup

//EtcdClient 客户端
type EtcdClient struct {
	config *config.EtcdConfig
	client *etcd_client.Client
	kv     *list.List
	mtx    sync.Mutex
	status int
}

//NewEtcdClient 创建etcd客户端
func NewEtcdClient(config *config.EtcdConfig) (*EtcdClient, error) {
	if config == nil {
		return nil, errors.New("NewEtcdClient: EtcdConfig invalid")
	}
	wg = &sync.WaitGroup{}
	etcdPodChan = make(chan *EtcdPod, config.Etcd2KakfaChanSize)
	return &EtcdClient{client: nil, config: config, status: EtcdInitial}, nil
}

//Close 关闭etcd客户端
func (cli *EtcdClient) close() {
	if cli.client != nil {
		cli.client.Close()
	}
}

func (cli *EtcdClient) open() error {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints:   cli.config.EndPoints,
		DialTimeout: time.Duration(cli.config.DialTimeOut) * time.Second,
	})

	if err != nil {
		logs.Error("InitEtcd: ", err)
		return err
	}
	cli.client = client
	return nil
}

func (cli *EtcdClient) processEtcdKV(key, value, kvt string) {
	switch kvt {
	case "PUT":
		cli.mtx.Lock()
		cli.kv.PushBack(&EtcdPod{key, value, kvt})
		cli.mtx.Unlock()
	case "DELETE":
		cli.mtx.Lock()
		cli.kv.PushBack(&EtcdPod{key, value, kvt})
		cli.mtx.Unlock()
	default:
		logs.Warn("EtcdClient::prcessEtcdKV: unkown operation: ", kvt, "[", key, "]=", value)
	}
}

func (cli *EtcdClient) transferToKafka() {
	for {
		for i := 0; i < cli.kv.Len(); i++ {
			e := cli.kv.Front()
			//transfer
			etcdPodChan <- e.Value.(*EtcdPod)
			//remove
			cli.mtx.Lock()
			cli.kv.Remove(e)
			cli.mtx.Unlock()
		}
		if cli.status == EtcdShutDown {
			break
		}
	}
}

func (cli *EtcdClient) isStatusKey(key string) bool {
	statusKey := cli.config.WatchKeyPrefix
	if !strings.HasPrefix(cli.config.WatchKeyPrefix, "/") {
		statusKey += "/"
	}
	all := statusKey + cli.config.StatusKey
	return all == key
}

func (cli *EtcdClient) load(timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	resp, err := cli.client.Get(ctx, cli.config.WatchKeyPrefix, etcd_client.WithPrefix())
	cancel()
	if err != nil {
		err = fmt.Errorf("EtcdClient::load: %s", err)
		return err
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		if !cli.isStatusKey(string(ev.Key)) {
			cli.processEtcdKV(string(ev.Key), string(ev.Value), "PUT")
		}
	}
	return nil
}

func (cli *EtcdClient) watcher() {
	defer wg.Done()
	//监测更新服务
	logs.Debug("EtcdClient::watcher: Watch start ....")
	err := cli.load(5)
	if err != nil {
		logs.Error("EtcdClient::watcher: ", err)
		return
	}
	logs.Debug("EtcdClient::watcher: Load complete ...")
	for {
		if cli.status == EtcdShutDown {
			return
		}

		rch := cli.client.Watch(context.Background(), cli.config.WatchKeyPrefix, etcd_client.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				// fmt.Printf("%s: %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if !cli.isStatusKey(string(ev.Kv.Key)) {
					cli.processEtcdKV(string(ev.Kv.Key), string(ev.Kv.Value), ev.Type.String())
				} else {
					value := string2Status(string(ev.Kv.Value))
					if value == EtcdShutDown {
						return
					} else if value == EtcdStop {
						for {
							if cli.status != EtcdStop {
								break
							}
							time.Sleep(time.Second)
						}
					}
				}
			}
		}
	}
}

func status2String(value int) string {
	switch value {
	case EtcdInitial:
		return "EtcdInitial"
	case EtcdStart:
		return "EtcdStart"
	case EtcdStop:
		return "EtcdStop"
	case EtcdShutDown:
		return "EtcdShutDown"
	}
	return "none"
}

func string2Status(value string) int {
	switch value {
	case "EtcdInitial":
		return EtcdInitial
	case "EtcdStart":
		return EtcdStart
	case "EtcdStop":
		return EtcdStop
	case "EtcdShutDown":
		return EtcdShutDown
	}
	return EtcdInitial
}

func (cli *EtcdClient) put(prefix, key, value string, timeout int) (err error) {
	if len(key) == 0 {
		err = errors.New("EtcdClient::put: key is empty")
		logs.Error(err)
		return
	}

	if len(prefix) != 0 && !strings.HasPrefix(prefix, "/") {
		prefix += "/"
	}
	all := prefix + key

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	_, err = cli.client.Put(ctx, all, value)
	cancel()
	if err != nil {
		err = fmt.Errorf("EtcdClient::put: %s", err)
	}
	return
}

//Start 启动Etcd监测服务
func (cli *EtcdClient) Start() {
	if cli.status == EtcdStart {
		logs.Info("EtcdClient::Start: Etcd Client already started")
		return
	} else if cli.status == EtcdStop {
		cli.status = EtcdStart
		return
	}
	err := cli.open()
	if err != nil {
		return
	}
	cli.status = EtcdStart
	cli.put(cli.config.WatchKeyPrefix, cli.config.StatusKey, status2String(cli.status), 1)
	cli.kv = list.New()
	wg.Add(1)
	go cli.watcher()
	go cli.transferToKafka()

	//Test API
	// msg := Msg{Topic: "logagent", Value: "./test.log"}
	// value, err := json.Marshal(msg)
	// if err != nil {
	// 	logs.Error("EtcdClient::Start:", err)
	// }
	// cli.put(cli.config.WatchKeyPrefix, "log", string(value), 1)
}

//Stop 暂停Etcd监测服务
func (cli *EtcdClient) Stop() {
	if cli.status == EtcdShutDown {
		logs.Info("EtcdClient::Start: Etcd Client already closed")
		return
	} else if cli.status == EtcdInitial {
		logs.Info("EtcdClient::Start: Etcd Client not started yet")
		return
	}
	cli.status = EtcdStop
	cli.put(cli.config.WatchKeyPrefix, cli.config.StatusKey, status2String(cli.status), 1)
}

//Restart 恢复Etcd监测服务
func (cli *EtcdClient) Restart() {
	if cli.status == EtcdStart {
		logs.Info("EtcdClient::Restart: Etcd Client already started")
	} else {
		logs.Info("EtcdClient::Restart:: Etcd Client start ...")
		cli.Start()
	}
}

//ShutDown 关闭Etcd监测服务
func (cli *EtcdClient) ShutDown() {
	if cli.status == EtcdShutDown {
		logs.Info("EtcdClient::ShutDown: EtcdClient already shutdown")
		return
	}
	cli.status = EtcdShutDown
	cli.put(cli.config.WatchKeyPrefix, cli.config.StatusKey, status2String(cli.status), 1)
	cli.kv = list.New()
	defer cli.close()
	//等待线程关闭
	wg.Wait()
}
