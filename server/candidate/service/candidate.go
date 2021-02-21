// @Author LiuYong
// @Created at 2021-02-05
// @Modified at 2021-02-05
package service

import (
	"candidate/tools/setting"
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
	cs.orderOp, err = orderRPCClient.NewClient()
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

	// TODO 创建候补存入缓存, 支付完成后, 支付服务通过队列通知候补服务, 候补将该候补信息写入mysql
	return resp.OrderOutsideID, err
}
