package worker

import (
	"context"
	"fmt"
	"time"
)

type Worker struct {
	worker chan struct{}
}

func New(capacity int) *Worker {
	return &Worker{worker: make(chan struct{}, capacity)}
}

func (w *Worker) Run(run func()) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	select {
	case w.worker <- struct{}{}:
		// 目前context还未使用
		ch := make(chan struct{}, 0)
		defer cancel()
		go func(ctx context.Context) {
			defer func() {
				<-w.worker
			}()
			run()

		}(ctx)
		select {
		case <-ch:
			fmt.Println("处理结束")
		case <-ctx.Done():
			fmt.Println("超时退出")
		}

	default:
		run()
	}
}
