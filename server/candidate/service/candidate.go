// @Author LiuYong
// @Created at 2021-02-05
// @Modified at 2021-02-05
package service

import (
	"candidate/tools/setting"
	"common/tools/logging"
	"rpc/pay/client/orderRPCClient"
	"rpc/pay/proto/orderRPCpb"
	"strconv"
	"time"
)

type candidateService struct {
	orderOp *orderRPCClient.OrderRPCClient
}

func NewCandidateService() (*candidateService, error) {
	cs := &candidateService{}
	var err error
	cs.orderOp, err = orderRPCClient.NewClientWithTargetAndMQHost(setting.RPCTarget.Order, setting.Kafka.Host)
	cs.orderOp.SetDealPayOK(payOK)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// CacheCandidate 创建候补订单存入缓存,后将返回点单号给前端, 前端根据订单号支付
func (c candidateService) CacheCandidate(userID, trainId uint, date string, passengers []string) (string, error) {
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

	// candidate := model.Candidate{}

	// // 先存储到redis
	// cc := cache.CandidateCache{}
	// cache2.Set(cc.GetKeyByOrderIDUnPay())
	return resp.OrderOutsideID, err
}

// payOK 用户支付完成有进行相关的处理
func payOK(payOKInfo *orderRPCClient.PayOKOrderInfo) {
	// 主要就是把缓存中的订单写入mysql
	logging.Info(payOKInfo.AffairID)
}
