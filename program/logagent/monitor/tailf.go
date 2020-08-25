package monitor

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

const (
	TailsOn  = 0
	TailsOff = 1
)

type TailfManager struct {
	tails  map[string]*TailInfo
	cancel context.CancelFunc
	status int
}

type TailInfo struct {
	filename string
	tailf    *tail.Tail
	cancel   context.CancelFunc
}

//NewTailfManager 创建Tail日志组件管理器
func NewTailfManager() *TailfManager {
	return &TailfManager{tails: make(map[string]*TailInfo), status: TailsOff}
}

func (m *TailfManager) create(filename string) (*tail.Tail, error) {
	if len(filename) == 0 {
		return nil, errors.New("TailfManager::create: filename is empty")
	}
	tailf, err := tail.TailFile(
		filename,
		tail.Config{
			ReOpen:    true,
			Follow:    true,
			Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
			MustExist: false,
			Poll:      true,
		})
	if err != nil {
		return nil, fmt.Errorf("TailfManager::create: %s", err)
	}
	return tailf, nil
}

func (m *TailfManager) close(key string) {

	if _, ok := m.tails[key]; ok {
		tailf := m.tails[key]
		tailf.cancel()
		delete(m.tails, key)
	}
}

func (m *TailfManager) monitor(ctx context.Context, key string) {

	tailf := m.tails[key].tailf
	for {
		select {
		case <-ctx.Done():
			logs.Info("TailfManager::monitor[", key, "] exit.")
			return
		case msg, ok := <-tailf.Lines:
			if !ok {
				logs.Error("TailfManager::monitor[", key, "]=", tailf.Filename, " close reopen")
				time.Sleep(time.Second)
			}
			//发送kafka
			fmt.Println("msg: ", msg)
		}
	}
}

func (m *TailfManager) processTailKv(key, value, kvt string) {
	switch kvt {
	case "PUT":
		tailf, err := m.create(value)
		if err != nil {
			logs.Error("TailfManager::processTailKv: ", err)
			return
		}
		m.close(key)
		ctx, cancel := context.WithCancel(context.Background())
		m.tails[key] = &TailInfo{value, tailf, cancel}
		go m.monitor(ctx, key)
	case "DELETE":
		m.close(key)
	default:
		logs.Error("TailfManager::processTailKv: Not Support Action ", kvt, " [", key, "]= ", value)
	}
}

//tailf 监测器
func (m *TailfManager) watcher(ctx context.Context) {
	for {
		select {
		case e, ok := <-etcdToTailChan:
			if ok {
				m.processTailKv(e.key, e.value, e.kvt)
			}
		case <-ctx.Done():
			logs.Info("TailfManager::watcher: exit")
			return
		}
	}
}

//Start 启动管理器
func (m *TailfManager) Start() {
	if m.status == TailsOn {
		logs.Info("TailfManager already started")
		return
	}
	m.status = TailsOn
	logs.Info("TailfManager starting ...")
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	go m.watcher(ctx)
}

//ShutDown 关闭管理器
func (m *TailfManager) ShutDown() {
	if m.status == TailsOff {
		logs.Info("TailfManager already shutdown")
		return
	}
	logs.Info("TailfManager shutting ...")
	for key := range m.tails {
		fmt.Printf("key :%s\n", key)
		m.close(key)
	}
	m.cancel()
	time.Sleep(time.Second)
	m = &TailfManager{tails: make(map[string]*TailInfo)}

}
