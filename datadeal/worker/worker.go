package worker

import (
	"context"
	"deal_data/service/mysql"
	"fmt"
	"go.uber.org/zap"
	"time"
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
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
		defer cancel()

		go func(ctx context.Context, news mysql.News) {
			ch := make(chan struct{}, 0)
			defer func() {
				<-w.worker
			}()
			go func(news mysql.News) {
				run(news)
				ch <- struct{}{}
			}(news)

			select {
			case <-ch:
			case <-ctx.Done():
				zap.L().Warn(fmt.Sprintf("图片下载超时退出,图片url:%s,新闻id:%d", news.Url, news.Id))
			}
		}(ctx, data)
	default:
		run(data)
	}
}
