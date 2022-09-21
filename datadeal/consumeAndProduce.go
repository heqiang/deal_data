package datadeal

import (
	"context"
	"deal_data/global"
	"deal_data/mysqlservice"
	"fmt"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
	"time"
)

var (
	Cond sync.Cond
)

func Produce(out chan<- mysqlservice.News) {
	for {
		Cond.L.Lock()
		for len(out) == 3 {
			Cond.Wait()
		}
		news, err := global.Db.Select()
		if err != nil {
			fmt.Println(err)
		}
		for _, data := range news {
			out <- data
		}

		Cond.L.Unlock()
		Cond.Signal()
		time.Sleep(time.Second)

	}
}

func Consume(in <-chan mysqlservice.News, ctx context.Context) {
	for {
		Cond.L.Lock()
		for len(in) == 0 {
			Cond.Wait()
		}
		data := <-in

		err := global.Db.UpdateNew(data.Id, 1)
		fmt.Println(fmt.Sprintf("正在处理数据:%d", data.Id))
		if err != nil {
			fmt.Println(err)
		}
		//TODO 协程超时退出
		go func(news mysqlservice.News) {
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error(fmt.Sprintf("数据处理异常:%s,err:%d,新闻id:%d,", err, debug.Stack(), news.Id))
					return
				} else {
					err = global.Db.UpdateNew(news.Id, 2)
					if err != nil {
						zap.L().Error(fmt.Sprintf("新闻处理状态更新2失败,err:%s,debugStack:%s,新闻id:%d", err, debug.Stack(), news.Id))
						return
					}
					fmt.Println(fmt.Sprintf("%d,数据处理结束", data.Id))
				}
			}()
			//数据处理的主要逻辑
			deal := NewDataDeal(global.Proxy, news.Direction)
			deal.TransNewsToJson(news)
			deal.download(news.Content, news.Id)

		}(data)

		Cond.L.Unlock()
		Cond.Signal()
		time.Sleep(time.Millisecond * 500)
	}
}
