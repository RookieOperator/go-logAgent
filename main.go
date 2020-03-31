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
	"code.rookieops.com/coolops/logAgent/etcd"
	"code.rookieops.com/coolops/logAgent/kafka"
	"code.rookieops.com/coolops/logAgent/tailLog"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
)

// 定义一个加载配置文件变量
var (
	cfg = new(config.AppConfig)
	err error
)

// 程序入口
func main() {
	// 0、加载配置文件
	err = ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		panic(err)
	}
	fmt.Println("load config success")
	// 1、初始化kafka
	err = kafka.InitKafka([]string{cfg.KafkaConfig.Address}, cfg.KafkaConfig.MaxSize)
	if err != nil {
		err = errors.New("init kafka failed")
		panic(err)
	}
	fmt.Println("init kafka success")

	// 2、初始化etcd
	err = etcd.InitEtcd([]string{cfg.EtcdConfig.Address})
	if err != nil {
		panic(err)
	}
	fmt.Println("init etcd success")
	// 2.1、etcd初始化成功后，从etcd中获取配置信息
	logEntry, err := etcd.GetConfFromEtcd(cfg.EtcdConfig.Topic)
	if err != nil {
		fmt.Println("get config from etcd failed. err:", err)
	}
	// 打印获取到的配置文件信息
	for addr, topic := range logEntry {
		fmt.Printf("addr: %v topic: %v\n", addr, topic)
	}

	// 3、初始化tailLog
	err = tailLog.Init(logEntry)
	if err != nil {
		fmt.Println("init tailLog failed. err:", err)
		panic(err)
	}
	fmt.Println("init tailLog success")
	// 4、开启一个goroutine实时去监听etcd中配置文件的变化
	newConfChan := tailLog.GetNewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	etcd.WatchConfFromEtcd(cfg.EtcdConfig.Topic, newConfChan)
	wg.Wait()
	//run()
}
