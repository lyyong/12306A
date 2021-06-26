// Package machine
// @Author LiuYong
// @Created at 2021-06-17
package machine

import (
	"candidate/model"
	"candidate/service/cache"
	"common/tools/logging"
	"context"
	cache2 "pay/tools/cache"
	"rpc/ticket/Client"
	"rpc/ticket/proto/ticketRPC"
	"time"
)

var tClient *client.TicketRPCClient

// SetupByDuration 开启定时抢票的功能, 时间间隔, 最好是小时为时间间隔
func SetupByDuration(ctx context.Context, d time.Duration, ticketURL string) {
	var err error
	tClient, err = client.NewClientWithTarget(ticketURL)
	if err != nil {
		logging.Error(err)
		return
	}
	ticker := time.NewTicker(d)
	go func() {
		defer func() {
			ticker.Stop()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tryGetTicketsByDatabase()
			}
		}
	}()
}

// SetupByTime 开启定时抢票功能, 通过时间点
func SetupByTime() {

}

// getTrainIDs 获取有候补需求的车次id
// func getTrainIDs()  {
// 	onceGetTrainIDs.Do(func() {
// 		cc := cache.CandidateCache{}
// 		// 先看缓存中有没有
// 		if cache2.Exists(cc.GetTrainIDSCacheKey()) {
// 			result, err := cache2.RedisConn.SMembers(context.Background(),cc.GetTrainIDSCacheKey()).Result()
// 			if err != nil {
// 				logging.Error("缓存中获取车次id出错: ",err)
// 				return
// 			}
// 			trainIDs = make([]uint,len(result))
// 			for i:=range result {
// 				a,err := strconv.Atoi(result[i])
// 				if err != nil {
// 					logging.Error("车次id由string转换int出错: ",err,"  ",result[i])
// 					return
// 				}
// 				trainIDs[i] = uint(a)
// 			}
// 			return
// 		}
// 		// 从数据库获取
// 		database.Client().Raw("select distinct train_id from candidates").Scan(&trainIDs)
// 		// 存入缓存
// 		err := cache2.RedisConn.SAdd(context.Background(), cc.GetTrainIDSCacheKey(), trainIDs).Err()
// 		if err != nil {
// 			logging.Error("候补车次id存入缓冲出错: ",err)
// 			return
// 		}
// 	})
// }

// tryGetTicketsByDatabase 尝试获取车票
func tryGetTicketsByDatabase() {
	// 当前日期
	thisDate := time.Now()
	// 获取订单
	orderIDs := model.GetCandidatesOrderIDs()
	for i := range orderIDs {
		candidates, err := model.GetCandidatesByOrderID(orderIDs[i].OrderID)
		if err != nil || len(candidates) == 0 {
			logging.Error("获取候补订单出错: ", err, "订单编号: ", orderIDs[i])
			continue
		}
		// 检查候补订单的时间
		if candidates[0].Date.Before(thisDate) || candidates[0].ExpireDate.Before(thisDate) {
			err := model.UpdateCandidatesState(orderIDs[i].OrderID, model.CandidateFail)
			if err != nil {
				logging.Error("修改候补状态出错: ", err, "订单编号: ", orderIDs[i])
				return
			}
		}
		// 获取票
		passengers := make([]*ticketRPC.Passenger, len(candidates))
		for j := range passengers {
			passengers[j] = &ticketRPC.Passenger{
				PassengerId:   uint32(candidates[j].PassengerID),
				PassengerName: candidates[j].PassengerName,
				SeatTypeId:    uint32(candidates[j].SeatTypeID),
			}
		}
		req := &ticketRPC.BuyTicketsRequest{
			TrainId:        uint32(candidates[0].TrainID),
			StartStationId: uint32(candidates[0].StartStationID),
			DestStationId:  uint32(candidates[0].DestStationID),
			Date:           candidates[0].Date.Format("2006-01-02"),
			Passengers:     passengers,
			OrderOuterId:   candidates[0].OrderID,
			UserId:         uint32(candidates[0].UserID),
		}
		tickets, err := tClient.BuyTickets(req)
		if err != nil {
			logging.Error("候补订单获取票出错: ", err, " 订单编号: ", candidates[0].OrderID)
			continue
		}
		// 存入数据库
		for i := range tickets.Response {
			for j := range candidates {
				if uint32(candidates[j].PassengerID) == tickets.Response[i].PassengerId {
					candidates[j].TicketID = uint(tickets.Response[i].TicketId)
					err := model.UpdateCandidate(candidates[i])
					if err != nil {
						logging.Error("更新订单信息出错: ", err, " 订单编号: ", orderIDs[i], " 乘客id: ", candidates[j].PassengerID)
						return
					}
				}
			}
		}
	}
}

// tryGetTicketsByCache 尝试获取通过缓存中的信息车票
func tryGetTicketsByCache() {
	// 今天开始往后30天
	thisDay := time.Now()
	cc := cache.CandidateCache{}
	// 获取订单
	trainIDs := model.GetCandidateTrainIDs()

	for i := range trainIDs {
		for j := 1; j <= 30; j++ {
			// redis 中通过车次和日期作为key存储一个链表, 一个链表节点就是一个组订单
			key := cc.GetKeyByTrainIDAndDate(trainIDs[i], thisDay.Add(time.Duration(j)*24*time.Hour).Format("2006-01-02"))
			candidates := make([]model.Candidate, 0)
			if cache2.Exists(key) {
				lLen, _ := cache2.RedisConn.LLen(context.Background(), key).Result()
				for k := 0; k < int(lLen); k++ {
					err := cache2.LIndex(key, k, &candidates)
					if err != nil {
						logging.Error(err)
						return
					}
					// 开始抢票
					passenger := make([]*ticketRPC.Passenger, len(candidates))
					for i := range passenger {
						passenger[i] = &ticketRPC.Passenger{
							PassengerId:   uint32(candidates[i].PassengerID),
							PassengerName: candidates[i].PassengerName,
							SeatTypeId:    uint32(candidates[i].SeatTypeID),
						}
					}
					req := &ticketRPC.BuyTicketsRequest{
						TrainId:        uint32(candidates[0].TrainID),
						StartStationId: uint32(candidates[0].StartStationID),
						DestStationId:  uint32(candidates[0].DestStationID),
						Date:           candidates[0].Date.Format("2006-01-02"),
						Passengers:     passenger,
						OrderOuterId:   candidates[0].OrderID,
						UserId:         uint32(candidates[0].UserID),
					}
					tickets, err := tClient.BuyTickets(req)
					if err != nil || len(tickets.Response) != len(candidates) {
						logging.Error("候补订单获取票出错: ", err, " 订单编号: ", candidates[0].OrderID)
						return
					}
					// 存入数据库
					for i := range tickets.Response {
						for j := range candidates {
							if uint32(candidates[j].PassengerID) == tickets.Response[i].PassengerId {
								candidates[j].TicketID = uint(tickets.Response[i].TicketId)
							}
						}
					}
					err = model.AddCandidates(candidates)
					if err != nil {
						logging.Error("有票的候补订单存入数据库出错: ", err)
						return
					}
					// 删除缓存
					err = cache2.LRem(key, k, &candidates)
					if err != nil {
						logging.Error("删除有票的候补订单出错: ", err)
						return
					}
				}
			} else {
				// 从数据库获取

			}

		}
	}
}
