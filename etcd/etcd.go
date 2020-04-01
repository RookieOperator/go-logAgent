package etcd

import (
	"code.rookieops.com/coolops/logAgent/logger"
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

// 定义一个结构体，用来保存获取到的信息
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// 声明一个全局的client变量，负责初始化成功后的etcd的一系列操作
var cli *clientv3.Client

// InitEtcd etcd的初始化操作
func InitEtcd(address string) (err error) {
	var etcdAddress = make([]string, 0)
	if strings.Contains(address, ";") {
		etcdAddress = strings.Split(address, ";")
	} else {
		etcdAddress = append(etcdAddress, address)
	}

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdAddress,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		logger.SugarLogger.Errorf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

// GetConfFromEtcd 从etcd中获取配置信息
func GetConfFromEtcd(key string) (logEntry []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var resp *clientv3.GetResponse
	resp, err = cli.Get(ctx, key)
	cancel()
	if err != nil {
		logger.SugarLogger.Errorf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		// 对取取到的配置文件进行反序列话
		err = json.Unmarshal(ev.Value, &logEntry)
		if err != nil {
			logger.SugarLogger.Errorf("Unmarshal failed for etcd value. err:%s", err)
			return
		}
	}
	return
}

// WatchConfFromEtcd 监控Etcd中配置文件的变化
func WatchConfFromEtcd(key string, newConfChan chan<- []*LogEntry) {
	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// 如果配置文件变了，就通知tailLog
			// 先解析新的配置文件
			var newConf []*LogEntry
			// 判断一下ev.Type，如果是delete操作，就向里面放一个空指针
			if ev.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					logger.SugarLogger.Errorf("new conf unmarshal failed, err:%s", err)
					continue
				}
			}
			// 将序列化好的数据传给channel
			newConfChan <- newConf
		}
	}
}
