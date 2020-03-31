package tailLog

import (
	"code.rookieops.com/coolops/logAgent/etcd"
	"fmt"
	"time"
)

var tailMgr *tailTaskMgr

// 管理任务的结构体
type tailTaskMgr struct {
	// 保存每次任务
	logEntry []*etcd.LogEntry
	// 将变化的配置文件保存到channel中
	newConfChan chan []*etcd.LogEntry
	// 定义一个map，用来存放
	taskMap map[string]*TailTask


}

// 每创建一个任务，就初始化一下，方便 热更新
func Init(logEntry []*etcd.LogEntry) (err error) {
	// 初始化tailMgr

	tailMgr = &tailTaskMgr{
		logEntry: logEntry,
		// 初始化为无缓冲区channel
		newConfChan: make(chan []*etcd.LogEntry),
		taskMap:     make(map[string]*TailTask, 16),
	}
	for _, logTask := range logEntry {
		// 循环初始化
		tailObj, err := NewTailTask(logTask.Path, logTask.Topic)
		if err != nil {
			fmt.Println(err)
		}
		// 将每次的taiObj都存入map，将path和topic组合作为键值
		mKey := fmt.Sprintf("%s_%s", logTask.Path, logTask.Topic)
		tailMgr.taskMap[mKey] = tailObj
	}

	// 开启监听配置文件变化goroutine
	go tailMgr.run()
	return
}

// 新增一个run方法来时刻监听channel变化，如果变化了，取channel中的值进行操作
func (t tailTaskMgr) run() {
	// 循环监听
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("get new conf ...", newConf)
			// 拿到新的配置文件，对配置文件进行判断
			for _, conf := range newConf {
				mKey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.taskMap[mKey]
				if !ok {
					tailObj, err := NewTailTask(conf.Path, conf.Topic)
					if err != nil {
						return
					}
					// 将每次的taiObj都存入map
					mKey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
					t.taskMap[mKey] = tailObj
					fmt.Printf("start new task %s_%s\n", conf.Path,conf.Topic)
				}
			}
			// 对比新旧配置文件，结束已经删除配置文件的goroutine
			for _, c1 := range t.logEntry {
				isDelete := true
				//fmt.Printf("删除task操作：c1.path: %v c1.topic: %v\n",c1.Path,c1.Topic)
				for _, c2 := range newConf {
					//fmt.Printf("删除task操作：c2.path: %v c2.topic: %v\n",c2.Path,c2.Topic)
					if c1.Path == c2.Path && c1.Topic == c2.Topic {
						isDelete = false
						continue
					}
				}
				// 如果不相等就删除goroutine，这里需要使用context
				fmt.Println("isDelete: ",isDelete)
				if isDelete {
					mKey := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					fmt.Printf("mKey: %s",mKey)
					t.taskMap[mKey].ctxCancelFunc()
				}
			}
		default:
			time.Sleep(time.Millisecond*50)
		}
	}
}

// 定义一个方法，使外部可以访问内部的私有变量
func GetNewConfChan() chan<- []*etcd.LogEntry {
	return tailMgr.newConfChan
}
