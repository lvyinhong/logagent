package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"llvvlv00.org/logagent/conf"
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

	//2、初始化kafka的链接
	fmt.Println(cfg.Address)
	if err := kafka.Init([]string{cfg.Address});err != nil {
		fmt.Println("初始化 kafka 失败 ", err)
		return
	}
	fmt.Println("初始化 kafka 成功!")

	//3、打开日志文件准备收集日志
	if err := taillog.Init(cfg.Filename);err!= nil {
		fmt.Println("初始化 taillog 失败 ", err)
		return
	}
	fmt.Println("初始化 taillog 成功!")

	// 4、读取日志发送到kafka
	run()
}


// 读取日志发送到kafka
func run() {
	//1、读取日志
	for {
		select {
		case line := <-taillog.ReadChan():
			// 2、发送到kafka
			kafka.SendToKafka(cfg.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}