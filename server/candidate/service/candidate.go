// Package service
// @Author LiuYong
// @Created at 2021-02-05
package service

import (
	"candidate/model"
	"candidate/service/cache"
	"candidate/tools/setting"
	"common/tools/logging"
	"errors"
	cache2 "pay/tools/cache"
	"rpc/pay/client/orderRPCClient"
	"rpc/pay/proto/orderRPCpb"
	"rpc/ticket/Client"
	"rpc/ticket/proto/ticketRPC"
	"rpc/user/userrpc"
	"strconv"
	"strings"
	"time"
)

type candidateService struct {
	orderOp  *orderRPCClient.OrderRPCClient
	ticketOp *client.TicketRPCClient
	userOp   *userrpc.Client
}

func NewCandidateService() (*candidateService, error) {
	cs := &candidateService{}
	var err error
	cs.orderOp, err = orderRPCClient.NewClientWithTargetAndMQHost(setting.RPCTarget.Order, setting.Kafka.Host)
	cs.orderOp.SetDealPayOK(payOK)
	if err != nil {
		return nil, err
	}
	cs.ticketOp, err = client.NewClientWithTarget(setting.RPCTarget.Ticket)
	if err != nil {
		return nil, err
	}
	cs.userOp = userrpc.NewClientWithTarget(setting.RPCTarget.User)
	return cs, nil
}

// CacheCandidate 创建候补订单存入缓存,后将返回点单号给前端, 前端根据订单号支付
func (c candidateService) CacheCandidate(userID, trainId, startStationID, destStationID uint, date, expire *time.Time, passengers []uint, seatTypeID int) (string, error) {
	// 检查是否有正在候补的订单
	cans, err := model.GetCandidates(map[string]interface{}{"user_id": userID, "state": model.CandidateIng})
	thisDay := time.Now()
	if err == nil && len(cans) != 0 && cans[0].ExpireDate.After(thisDay) {
		return "", errors.New("已近存在别的候补订单")
	}
	for i := range cans {
		if cans[i].ExpireDate.Before(thisDay) {
			cans[i].State = model.CandidateFail
			model.UpdateCandidate(cans[i])
		}
	}

	money := 100
	// 创建订单, 获得外部id
	resp, err := c.orderOp.Create(&orderRPCpb.CreateRequest{
		UserID:         uint64(userID),
		Money:          int64(money),
		AffairID:       "CAN" + time.Now().Format("2006-01-02-15:04:05.000-") + strconv.Itoa(int(userID)),
		ExpireDuration: 1800,
		CreatedBy:      setting.Server.Name,
	})
	if err != nil {
		return "", err
	}

	listPassenger, err := c.userOp.ListPassenger(userID)
	if err != nil {
		return "", err
	}
	allPassenger := make(map[uint]*userrpc.Passenger)
	for i := range listPassenger {
		allPassenger[listPassenger[i].Id] = listPassenger[i]
	}
	candidates := make([]*model.Candidate, len(passengers))
	for i := range candidates {
		if _, ok := allPassenger[passengers[i]]; !ok {
			return "", errors.New("乘客id出错")
		}
		candidates[i] = &model.Candidate{
			Date:              *date,
			TrainID:           trainId,
			OrderID:           resp.GetOrderOutsideID(),
			UserID:            userID,
			PassengerID:       passengers[i],
			PassengerName:     allPassenger[passengers[i]].Name,
			CertificateNumber: allPassenger[passengers[i]].CertificateNumber,
			StartStationID:    startStationID,
			DestStationID:     destStationID,
			ExpireDate:        *expire,
			SeatTypeID:        seatTypeID,
			State:             0,
		}
	}

	// 先存储到redis
	cc := cache.CandidateCache{}
	err = cache2.Set(cc.GetKeyByOrderIDUnPay(resp.OrderOutsideID), candidates, 1800)
	if err != nil {
		logging.Error("候补订单写入缓存出错: ", err)
		return "", err
	}
	return resp.OrderOutsideID, err
}

// ReadCandidate 获取到候补订单的信息
func (c candidateService) ReadCandidate(userID uint) []*model.Candidate {
	cans, err := model.GetCandidates(map[string]interface{}{"user_id": userID, "state": []int{model.CandidateIng, model.CandidateNotCash}})
	if err != nil {
		logging.Error("读取候补订单出错: ", err)
	}
	return cans
}

// Cash 兑现候补订单
func (c candidateService) Cash(ids []uint) error {
	// 修改票的状态
	cans, err := model.GetCandidates(map[string]interface{}{"id": ids})
	if err != nil {
		logging.Error(err)
		return err
	}
	ticketIDs := make([]uint32, len(cans))
	for i := range cans {
		ticketIDs[i] = uint32(cans[i].TicketID)
	}
	req := &ticketRPC.UpdateStateRequest{
		TicketsId: ticketIDs,
		State:     0, // 0为购票成功
	}
	_, err = c.ticketOp.UpdateTicketsState(req)
	if err != nil {
		return err
	}
	// 写回候补订单状态
	for i := range cans {
		cans[i].State = model.CandidateSuccess
		err := model.UpdateCandidate(cans[i])
		if err != nil {
			logging.Error(err)
			return err
		}
	}
	return nil
}

// payOK 用户支付完成有进行相关的处理
func payOK(payOKInfo *orderRPCClient.PayOKOrderInfo) {
	// 主要就是把缓存中的订单写入mysql
	if !strings.HasPrefix(payOKInfo.AffairID, "CAN") {
		return
	}
	logging.Info("收到支付完成的通知, 订单编号: ", payOKInfo.AffairID)
	cc := cache.CandidateCache{}
	candidates := make([]model.Candidate, 0)
	// 获取缓存中的候补订单
	err := cache2.Get(cc.GetKeyByOrderIDUnPay(payOKInfo.OutsideID), &candidates)
	if err != nil {
		logging.Error("获取缓存中的候补订单出错: ", err)
		return
	}
	// 存储到数据库中
	err = model.AddCandidates(candidates)
	if err != nil {
		logging.Error("候补存入数据库出错: ", err)
		return
	}
	// 删除缓存
	_, err = cache2.Delete(cc.GetKeyByOrderIDUnPay(payOKInfo.OutsideID))
	if err != nil {
		logging.Error("删除缓存中的候补订单出错: ", err)
		return
	}
	// 写入缓存中的链表
	// err = cache2.LPush(cc.GetKeyByTrainIDAndDate(candidates[0].TrainID, candidates[0].Date.Format("2006-01-02")), candidates)
	// if err != nil {
	// 	logging.Error("候补订单写入缓存链表出错: ", err)
	// 	return
	// }
	// 车次写入set中
	// cache2.RedisConn.SAdd(context.Background(),cc.GetTrainIDSCacheKey(),candidates[0].TrainID)
}
