package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

// 基于sarama第三方开发的kafka server
func main() {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer([]string{"192.168.10.10:9092"}, nil)
	if err != nil {
		fmt.Println("Failed to start consumer: %s", err)
		return
	}

	partitionList, err := consumer.Partitions("web_log")	//获得该topic所有分区
	if err != nil {
		fmt.Println("Failed to get the list of partition:", err)
		return
	}
	fmt.Println(partitionList)

	for partition := range partitionList {
		// 获得具体分区的consumer
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("Failed to start consumer for partition %d\n : %s\n", partition, err)
			return
		}
		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) { //为每个分区开一个go协程去取值
			for msg :=range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()
	if err := consumer.Close();err != nil {
		println(err)
		return
	}
}