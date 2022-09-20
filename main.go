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
	global.InitConfig("dev")
	err = log.InitLogger(global.ServerConfig.LogConfig, "dev")
	if err != nil {
		panic(fmt.Sprintf("日志初始化错误,err:%s", err))
	}
	global.Db = mysqlservice.NewMysqlConn(global.ServerConfig.MysqlConfig)

	quit := make(chan bool)
	newsChan := make(chan mysqlservice.News, 15)
	datadeal.Cond.L = new(sync.Mutex)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go datadeal.Produce(newsChan)
	for i := 0; i < 2; i++ {
		go datadeal.Consume(newsChan, ctx)
	}

	<-quit

}
