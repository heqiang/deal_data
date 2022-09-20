package datadeal

import (
	"deal_data/global"
	"deal_data/mysqlservice"
	"fmt"
	"go.uber.org/zap"
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

func Consume(in <-chan mysqlservice.News) {
	for {
		Cond.L.Lock()
		for len(in) == 0 {
			Cond.Wait()
		}
		data := <-in
		err := global.Db.UpdateNew(data.Id, 1)
		if err != nil {
			fmt.Println(err)
		}
		// 处理数据
		go func(news mysqlservice.News) {
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error(fmt.Sprintf("数据处理异常:%s,新闻id:%d", err, news.Id))
					return
				} else {
					err = global.Db.UpdateNew(news.Id, 2)
					if err != nil {
						zap.L().Error(fmt.Sprintf("新闻处理状态更新2失败,err:%s,新闻id:%d", err, news.Id))
						return
					}
				}
			}()
			deal := NewDataDeal("", news.Direction)
			deal.download(news.Content, news.Id)
			deal.TransNewsToJson(news)
		}(data)

		Cond.L.Unlock()
		Cond.Signal()
		time.Sleep(time.Millisecond * 500)
	}
}
