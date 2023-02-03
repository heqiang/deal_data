package worker

import (
	"deal_data/service/mysql"
)

type Worker struct {
	worker chan struct{}
}

func New(capacity int) *Worker {
	return &Worker{worker: make(chan struct{}, capacity)}
}

func (w *Worker) Run(run func(mysql.News), data mysql.News) {
	select {
	case w.worker <- struct{}{}:
		// 目前context还未使用
		go func() {
			defer func() {
				<-w.worker
			}()
			run(data)
		}()
	default:
		run(data)
	}
}
