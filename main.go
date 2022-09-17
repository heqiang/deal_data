package main

import (
	"deal_data/datadeal"
	"deal_data/global"
	"deal_data/mysqlservice"
	"sync"
)

func main() {
	global.InitConfig("dev")
	global.Db = mysqlservice.NewMysqlConn(global.ServerConfig.MysqlConfig)

	quit := make(chan bool)
	newsChan := make(chan mysqlservice.News, 15)
	datadeal.Cond.L = new(sync.Mutex)

	go datadeal.Produce(newsChan)
	for i := 0; i < 10; i++ {
		go datadeal.Consume(newsChan)
	}

	<-quit
}
