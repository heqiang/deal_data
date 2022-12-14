package main

import (
	"context"
	"deal_data/datadeal"
	"deal_data/global"
	"deal_data/log"
	"deal_data/mysqlservice"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		msg := fmt.Sprintf("根目录路径获取错误:%s", err)
		zap.L().Error(msg)
		return
	}
	global.AbsDataPath = filepath.Join(rootPath, "data")
	// 初始化配置
	global.InitConfig("dev")
	err = log.InitLogger(global.ServerConfig.LogConfig, "dev")
	if err != nil {
		panic(fmt.Sprintf("日志初始化错误,err:%s", err))
	}
	global.Db = mysqlservice.NewMysqlConn(global.ServerConfig.MysqlConfig)

	quit := make(chan bool)
	newsChan := make(chan mysqlservice.News, 15)

	global.Cond.L = new(sync.Mutex)
	// 目前context还未使用
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//生产者消费者模式去处理数据
	pipeline := datadeal.NewPipeline()
	go pipeline.Produce(newsChan)
	pipeline.Consume(newsChan, ctx)

	<-quit

}
