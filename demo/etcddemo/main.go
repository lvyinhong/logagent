package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// etcd 使用put命令设置键值对数据，get命令用来根据key获取值
func main(){
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"192.168.10.10:2379"},
		Username:"test",
		Password:"123456",
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		//handle error!
		fmt.Printf("connect to etcd failed, err :%v\n", err)
		return
	}

	defer cli.Close()

	//put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_,err = cli.Put(ctx, "name","llvvlv00ss")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed， err: %v\n", err)
		return
	}
	fmt.Println("put success!")

	//get
	//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//resp, err := cli.Get(ctx, "name")
	//cancel()
	//if err != nil {
	//	fmt.Printf("get from etcd failed， err: %v\n", err)
	//	return
	//}
	//
	//for _,ev := range resp.Kvs {
	//	fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	//}
	//
	//fmt.Println("get success!")
}
