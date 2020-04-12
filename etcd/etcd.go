package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path string `json:"path"`
	Topic string `json:"topic"`
}

// 初始化ectd的函数
func Init(addrs string, timeout time.Duration)(err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:[]string{addrs},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v\n", err)
		return
	}
	return
}

// 从etcd中根据key获取配置项
func GetConf(key string) (logEntrys []*LogEntry,err error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		err := json.Unmarshal(ev.Value, &logEntrys)
		if err != nil {
			fmt.Printf("unmarshal etcd value failed, err: %v\n", err)
		}
	}
	return
}

func WatchConf(key string, newConfCh chan <- []*LogEntry) {
	ch := cli.Watch(context.Background(), key)
	//从通道中尝试读取值(监视的信息)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Type:%v key:%v value:%v", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			//通知 taillog.tskMgr
			var newConf []*LogEntry
			if evt.Type == clientv3.EventTypeDelete {
				// 如果是删除操作
			}
			err := json.Unmarshal(evt.Kv.Value, &newConf)
			if err!= nil{
				fmt.Printf("unmarshal failed, err:%v\n", err)
				continue
			}
			fmt.Printf("get new conf:%v\n", newConf)
			newConfCh <- newConf
		}
	}
}