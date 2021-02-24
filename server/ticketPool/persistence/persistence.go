// @Author: KongLingWen
// @Created at 2021/2/23
// @Modified at 2021/2/23

package persistence

import (
	"ticketPool/model"
)


var persistenceChan chan *PersistentRequest

func init(){
	persistenceChan = make(chan *PersistentRequest, 10000)
	doPersistence()
}

type PersistentRequest struct {
	Option     string
	TrainId    uint32
	Date       string
	SeatTypeId uint32
	Key        uint64
	Value      []string
}


func Do(req *PersistentRequest){
	persistenceChan <- req
}

func doPersistence() {
	go func() {
		for {
			req := <-persistenceChan
			switch req.Option {
			case "DELETE":
				seat := &model.Seat{
					Key:        0,
					TrainId:    req.TrainId,
					Date:       req.Date,
					SeatTypeId: req.SeatTypeId,
				}
				model.DeleteSeat(seat, req.Value)
			case "INSERT":
				seats := make([]model.Seat, len(req.Value))
				for i := 0; i < len(req.Value); i++ {
					seats[i] = model.Seat{
						Key:        req.Key,
						TrainId:    req.TrainId,
						Date:       req.Date,
						SeatTypeId: req.SeatTypeId,
						SeatInfo:   req.Value[i],
					}
				}
				model.InsertSeat(&seats)
			}
		}
	}()
}