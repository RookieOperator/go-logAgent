/*
# #########################################
# Author:                     joker
# Blog:                       https://coolops.cn
# Date:						  2020-03-27
# FileName:					  main.go
# Description:	              日志收集客户端
# Copyright ©:
# #########################################
*/
package main

import (
	"code.rookieops.com/coolops/logAgent/config"
	"code.rookieops.com/coolops/logAgent/kafka"
	"code.rookieops.com/coolops/logAgent/tailLog"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

func run(){
	// 1、不断的读取日志
	for{
		select {
		case data := <-tailLog.ReadLine():
			// 2、向kafka里发送日志
			kafka.SendMsgToKafka(cfg.KafkaConfig.Topic, data.Text)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}

}

// 定义一个加载配置文件变量
var (
	cfg = new(config.AppConfig)
	err error
)

// 程序入口
func main(){
	// 0、加载配置文件
	err = ini.MapTo(cfg,"./config/config.ini")
	if err != nil {
		panic(err)
	}
	fmt.Println("load config success")
	// 1、初始化kafka
	err=kafka.InitKafka([]string{cfg.KafkaConfig.Address})
	if err != nil{
		err = errors.New("init kafka failed")
		panic(err)
	}
	fmt.Println("init kafka success")
	// 2、初始化tailLog
	err = tailLog.InitTailLog(cfg.TailLogConfig.FileName)
	if err != nil {
		err = errors.New("init tailLog failed")
		panic(err)
	}
	fmt.Println("init tailLog success")
	// 3、程序启动
	run()
}
