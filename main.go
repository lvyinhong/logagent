package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"llvvlv00.org/logagent/conf"
	"llvvlv00.org/logagent/etcd"
	"llvvlv00.org/logagent/kafka"
	"llvvlv00.org/logagent/taillog"
	"llvvlv00.org/logagent/utils"
	"sync"
	"time"
)

var (
	cfg = new(conf.AppConf)	//配置文件句柄
)

//logagent程序入口
func main() {

	//1、加载配置文件
	if err := ini.MapTo(cfg, "./conf/config.ini"); err != nil {
		fmt.Println("加载配置文件失败 ", err)
		return
	}
	fmt.Println("加载配置文件成功!")

	//2、初始化kafka的链接
	if err := kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.MaxSize);err != nil {
		fmt.Println("初始化 kafka 失败 ", err)
		return
	}
	fmt.Println("初始化 kafka 成功!")

	//3、初始化etcd
	if err := etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout) * time.Second); err != nil {
		fmt.Println("初始化 etcd 失败 ", err)
		return
	}
	fmt.Println("初始化 etcd 成功!")

	//3.1 从etcd中获取日志收集项的配置信息
	// 实现每一个logagent都拉去自己独有的配置，所有根据ip地址做区分
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipStr);

	logEntrys,err := etcd.GetConf(etcdConfKey)
	if err!= nil {
		fmt.Printf("etcd.GetConf failed, err:%v\n", err)
		return
	}
	if len(logEntrys) == 0 {
		fmt.Println("not found logEntrys conf")
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntrys)
	//4 收集日志发往Kafka
	taillog.Init(logEntrys)

	//3.2 派一个哨兵去监视日志收集项的变化(有变化及时通知logagent,实现配置热加载)
	newConfCh := taillog.NewConfCh()	//从taillog 包中获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		etcd.WatchConf(etcdConfKey,newConfCh) //哨兵发现最新的配置信息会通知通道
		wg.Done()
	}()

	wg.Wait()
}

