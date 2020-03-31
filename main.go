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
	"code.rookieops.com/coolops/logAgent/logger"
	"code.rookieops.com/coolops/logAgent/tailLog"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"sync"
)

// 定义一个加载配置文件变量
var (
	cfg = new(config.AppConfig)
	err error
	sugarLogger *zap.SugaredLogger
)

// 程序入口
func main() {
	// 0、加载配置文件
	err = ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		panic(err)
	}
	fmt.Println("load config success")

	// 1、初始化打印日志
	logger.InitLogger(cfg.LogConfig.Path, cfg.LogConfig.MaxSize, cfg.LogConfig.MaxBackups, cfg.LogConfig.MaxAge, cfg.LogConfig.Compress)

	// 2、初始化kafka
	err = kafka.InitKafka(cfg.KafkaConfig.Address, cfg.KafkaConfig.MaxSize)
	if err != nil {
		logger.SugarLogger.Errorf("init kafka failed. err:%s",err)
		panic(err)
	}
	logger.SugarLogger.Info("init kafka success")

	// 3、初始化etcd
	err = etcd.InitEtcd([]string{cfg.EtcdConfig.Address})
	if err != nil {
		panic(err)
	}
	logger.SugarLogger.Info("init etcd success")

	// 3.1、etcd初始化成功后，从etcd中获取配置信息
	logEntry, err := etcd.GetConfFromEtcd(cfg.EtcdConfig.Topic)
	if err != nil {
		logger.SugarLogger.Errorf("get config from etcd failed. err:%s", err)
		panic(err)
	}
	logger.SugarLogger.Info("get config from etcd success.")

	// 4、初始化tailLog
	err = tailLog.Init(logEntry)
	if err != nil {
		logger.SugarLogger.Errorf("init tailLog failed. err:%s", err)
		panic(err)
	}
	logger.SugarLogger.Info("init tailLog success")

	// 5、开启一个goroutine实时去监听etcd中配置文件的变化
	newConfChan := tailLog.GetNewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	etcd.WatchConfFromEtcd(cfg.EtcdConfig.Topic, newConfChan)
	wg.Wait()
}
