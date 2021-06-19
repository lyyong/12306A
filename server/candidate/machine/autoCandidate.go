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
	"pay/tools/database"
	"rpc/ticketPool/Client"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"time"
)

var tpClient *Client.TPRPCClient

// SetupByDuration 开启定时抢票的功能, 时间间隔, 最好是小时为时间间隔
func SetupByDuration(ctx context.Context, d time.Duration, ticketPoolURL string) {
	var err error
	tpClient, err = Client.NewClientWithTarget(ticketPoolURL)
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

// tryGetTickets 尝试获取车票
func tryGetTickets() {
	// 今天开始往后30天
	thisDay := time.Now()
	cc := cache.CandidateCache{}
	// 获取需要候补的车次
	var trainIDs []uint
	database.Client().Raw("select distinct train_id from candidates").Scan(&trainIDs)
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
					passengers := make([]*ticketPoolRPC.PassengerInfo, len(candidates))
					for h := range passengers {
						passengers[h] = &ticketPoolRPC.PassengerInfo{
							PassengerId:   uint32(candidates[h].PassengerID),
							PassengerName: "test",
							SeatTypeId:    1,
							ChooseSeat:    "A",
						}
					}
					req := &ticketPoolRPC.GetTicketRequest{
						TrainId:        uint32(candidates[0].TrainID),
						StartStationId: uint32(candidates[0].StartStationID),
						DestStationId:  uint32(candidates[0].DestStationID),
						Date:           candidates[0].Date.Format("2006-01-02"),
						Passengers:     passengers,
					}
					// TODO 换新的抢票RPC接口
					ticket, err := tpClient.GetTicket(req)
					if err != nil || len(ticket.GetTickets()) != len(candidates) {
						continue
					}
					ticket.GetTickets()
				}
			} else {
				// 从数据库获取

			}
		}
	}
}
