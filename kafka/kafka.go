package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// 专门往kafka写日志的模块

type logData struct {
	topic string
	data string
}

var (
	client sarama.SyncProducer	//声明一个全局的链接kafka的生产者
	logDataChan chan *logData
)

// 初始化client
func Init(addrs []string, maxSize int)(err error) {
	config := sarama.NewConfig()
	//tailf 包使用
	config.Producer.RequiredAcks = sarama.WaitForAll	//发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true // 成功交付的消息将在success channel返回

	// 链接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	// 初始化logDataChan
	logDataChan = make(chan *logData, maxSize)

	// 后台开启goroutine 从通道中取数据发往kafka
	go sendToKafka()
	return
}

// 把日志数据发送到日志channel中
func SendToChan(topic, data string)  {
	msg := &logData{
		topic:topic,
		data:data,
	}
	logDataChan <- msg
}

// 消费通道数据，往kafka发送日志数据
func sendToKafka() {
	for {
		select {
		case ld:=<-logDataChan:
			//构造一个消息
			msg := &sarama.ProducerMessage{
				Topic:ld.topic,
				Value:sarama.StringEncoder(ld.data),
			}
			// 发送到kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err: ",err)
				return
			}
			fmt.Printf("pid: %v offset:%v\n",  pid, offset)

		default:
			time.Sleep(time.Millisecond*50)
		}
	}
}