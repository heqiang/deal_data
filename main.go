package main

import (
	"deal_data/config"
	"deal_data/datadeal"
	"deal_data/log"
	mysqlservice2 "deal_data/service/mysql"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	rootPath, err := os.Getwd()
	if err != nil {
		msg := fmt.Sprintf("根目录路径获取错误:%s", err)
		zap.L().Error(msg)
		return
	}
	config.AbsDataPath = filepath.Join(rootPath, "data")
	// 初始化配置
	config.InitConfig("dev")
	err = log.InitLogger(config.Conf.LogConfig, "dev")
	if err != nil {
		panic(fmt.Sprintf("日志初始化错误,err:%s", err))
	}
	mysqlservice2.Db = mysqlservice2.NewMysqlConn(config.Conf.MysqlConfig)

	quit := make(chan bool)
	newsChan := make(chan mysqlservice2.News, 15)

	config.Cond.L = new(sync.Mutex)

	//生产者消费者模式去处理数据
	pipeline := datadeal.NewPipeline()
	go pipeline.Produce(newsChan)

	pipeline.Consume(newsChan)

	<-quit

}
