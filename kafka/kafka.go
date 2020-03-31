package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// 声明一个全局变量client，用来做消息发送
var (
	client   sarama.SyncProducer
	dataChan chan *logData
)

type logData struct {
	topic string
	data  string
}

func SendToChan(topic, data string) {
	// 构造数据
	logSendData := &logData{
		topic: topic,
		data:  data,
	}
	// 将数据发送到channel中
	dataChan <- logSendData
}

// InitKafka kafka初始化函数
func InitKafka(address []string, maxSize int) (err error) {
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
	// 初始化dataChan
	dataChan = make(chan *logData, maxSize)
	// 开启一个goroutine不断接受channel中的数据
	go sendMsgToKafka()
	return
}

// SendMsgToKafka 发送消息到kafka
func sendMsgToKafka() {
	// 从channel中取日志信息发送到kafka
	fmt.Println("开始发送消息到kafka")
	for {
		select {
		case logConf := <-dataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = logConf.topic
			msg.Value = sarama.StringEncoder(logConf.data)

			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
