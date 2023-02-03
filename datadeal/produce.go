package datadeal

import (
	"deal_data/config"
	"deal_data/service/mysql"
	"go.uber.org/zap"
	"time"
)

func (p *Pipeline) Produce(out chan<- mysql.News) {
	for {
		config.Cond.L.Lock()
		news, err := mysql.Db.SelectZero()

		if err != nil {
			zap.L().Error("数据处理错误")
			continue
		}

		if len(news) == 0 {
			time.Sleep(time.Second)
			continue
		}

		for _, data := range news {
			out <- data
		}
		config.Cond.L.Unlock()
	}
}
