package datadeal

import (
	"deal_data/datadeal/worker"
	"deal_data/global"
	"deal_data/mysqlservice"
	"fmt"
	"sync"
	"time"
)

var Cond sync.Cond

type Pipeline struct {
	w *worker.Worker
}

func NewPipeline() *Pipeline {
	return &Pipeline{w: worker.New(10000)}
}

func (p *Pipeline) Produce(out chan<- mysqlservice.News) {
	for {
		news, err := global.Db.Select()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		for _, data := range news {
			out <- data
		}
	}
}
