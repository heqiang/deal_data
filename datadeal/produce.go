package datadeal

import (
	"deal_data/datadeal/worker"
	"deal_data/global"
	"deal_data/mysqlservice"
	"fmt"
	"time"
)

type Pipeline struct {
	w *worker.Worker
}

func NewPipeline() *Pipeline {
	return &Pipeline{w: worker.New(10000)}
}

func (p *Pipeline) Produce(out chan<- mysqlservice.News) {
	for {
		global.Cond.L.Lock()
		news, err := global.Db.Select()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		if len(news) > 0 {
			mysqlservice.MaxId = news[len(news)-1].Id
		}

		for _, data := range news {
			out <- data
		}
		global.Cond.L.Unlock()
	}
}
