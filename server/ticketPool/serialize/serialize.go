// @Author: KongLingWen
// @Created at 2021/6/19
// @Modified at 2021/6/19

package serialize

import (
	"common/tools/logging"
	"encoding/json"
	"io/ioutil"
	"ticketPool/ticketpool"
	"time"
)

func Serialize() {
	go func() {
		for {
			tp := ticketpool.Tp
			tp.RWLock.Lock()
			logging.Info("开始序列化")
			st := time.Now()
			res, err := json.Marshal(&tp)
			if err != nil {
				logging.Error(err)
			}
			err = ioutil.WriteFile("TicketPoolData.json", res, 0644)
			if err != nil {
				logging.Error(err)
			}
			logging.Info("序列化结束, 耗时: ", time.Now().Sub(st).Milliseconds())
			tp.RWLock.Unlock()
			time.Sleep(5 * time.Second)
		}
	}()
}
