package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"llvvlv00.org/logagent/conf"
	"llvvlv00.org/logagent/etcd"
	"llvvlv00.org/logagent/kafka"
	"llvvlv00.org/logagent/taillog"
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
	logEntrys,err := etcd.GetConf("/xxx/")
	if err!= nil {
		fmt.Printf("etcd.GetConf failed, err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntrys)
	for index, value := range logEntrys {
		fmt.Printf("index:%v value:%v\n", index, value)
	}
	//3.2 派一个哨兵去监视日志收集项的变化(有变化及时通知logagent,实现配置热加载)

	//4 收集日志发往Kafka
	taillog.Init(logEntrys)
}

