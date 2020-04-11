package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"192.168.10.10:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		//handle error!
		fmt.Printf("connect to etcd failed, err :%v\n", err)
		return
	}

	defer cli.Close()

	// 建立一个监听 key=name 变动的通道
	rch := cli.Watch(context.Background(), "name")	//<-chan watchResponse

	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s key: %s value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}