package datadeal

import (
	"deal_data/global"
	"deal_data/mysqlservice"
	"fmt"
	"sync"
	"time"
)

var Cond sync.Cond

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
