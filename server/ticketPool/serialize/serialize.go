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
			res, err := json.Marshal(&tp)
			if err != nil {
				logging.Error(err)
			}
			err = ioutil.WriteFile("TicketPoolData.json", res, 0644)
			if err != nil {
				logging.Error(err)
			}
			time.Sleep(1000)
		}
	}()
}