package datadeal

import (
	"deal_data/global"
	"deal_data/service/mysql"
	"go.uber.org/zap"
	"time"
)

func (p *Pipeline) Produce(out chan<- mysql.News) {
	for {
		global.Cond.L.Lock()
		news, err := global.Db.SelectZero()

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
		global.Cond.L.Unlock()
	}
}
