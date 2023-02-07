package worker

import "deal_data/service/mysql"

type Worker struct {
	worker chan struct{}
}

func New(capacity int) *Worker {
	return &Worker{worker: make(chan struct{}, capacity)}
}

func (w *Worker) Run(run func(news mysql.News), data mysql.News) {
	select {
	case w.worker <- struct{}{}:
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
