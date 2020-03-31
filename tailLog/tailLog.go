package tailLog

import (
	"code.rookieops.com/coolops/logAgent/kafka"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
)

//var tailsObj *tail.Tail

//由于会有多个日志目录，多个任务，所以每次初始化的时候都会单独生成一个*tail.Tail
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail

	// 定义一个结束的context
	ctx       context.Context
	ctxCancelFunc context.CancelFunc
}

// NewTailTask 开启任务收集初始化操作
func NewTailTask(path, topic string) (tailObj *TailTask, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:  path,
		topic: topic,
		// 定义结束
		ctx:       ctx,
		ctxCancelFunc: cancel,
	}
	err = tailObj.init()
	if err != nil {
		fmt.Println("init failed, err:", err)
		return
	}
	return
}

// 创建一个init()方法
func (t TailTask) init() (err error) {
	t.instance, err = tail.TailFile(t.path, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	// 开启获取日志任务
	go t.run()
	return
}

// 创建一个run方法，用于获取日志发送到Kafka
func (t TailTask) run() {
	for {
		select {
		case <- t.ctx.Done():
			fmt.Printf("stop task %s_%s success\n",t.path,t.path)
			return
		case data := <-t.instance.Lines:
			// 2、为了加速处理，将收集信息发送到channel中
			//kafka.SendMsgToKafka(t.topic, data.Text)
			kafka.SendToChan(t.topic, data.Text)
		//default:
		//	time.Sleep(time.Millisecond * 500)
		}
	}
}
