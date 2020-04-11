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