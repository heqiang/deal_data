package datadeal

import (
	"deal_data/global"
	"deal_data/mysqlservice"
	"fmt"
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
		fmt.Println(data)

		Cond.L.Unlock()
		Cond.Signal()
		time.Sleep(time.Millisecond * 500)
	}
}
