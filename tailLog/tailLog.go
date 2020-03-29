package tailLog

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var tailsObj *tail.Tail

// InitTailLog TailLog初始化函数
func InitTailLog(filename string)(err error){
	tailsObj, err = tail.TailFile(filename, tail.Config{
		ReOpen: true,
		Follow: true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	return
}

func ReadLine() <-chan *tail.Line{
	return tailsObj.Lines
}