package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 基于sarama第三方开发的kafka client
func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll	//赋值为-1	： producer 在follower 副本确认接收到数据后才算一次发送完成
	config.Producer.Partitioner = sarama.NewRandomPartitioner //	写到随机分区中。默认8个分区
	config.Producer.Return.Successes = true

	// 构建一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = `nimabi`
	msg.Value = sarama.StringEncoder("this is a good test??????")

	// 链接kafka server
	client, err := sarama.NewSyncProducer([]string{"192.168.10.10:9092"}, config)
	if err != nil {
		fmt.Println("producer close err, ", err)
		return
	}

	defer client.Close()
	// 发送一个消息
	pid, offset, err := client.SendMessage(msg)

	if err != nil {
		fmt.Println("Send Message Failed, ", err)
		return
	}

	fmt.Printf("分区ID:%v, offset:%v \n", pid, offset)
}