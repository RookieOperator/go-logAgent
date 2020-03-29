package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 声明一个全局变量client，用来做消息发送
var client sarama.SyncProducer

// InitKafka kafka初始化函数
func InitKafka(address []string)(err error){
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	return
}

// SendMsgToKafka 发送消息到kafka
func SendMsgToKafka(topic,data string){
// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}