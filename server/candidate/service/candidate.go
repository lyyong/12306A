// @Author LiuYong
// @Created at 2021-02-05
// @Modified at 2021-02-05
package service

import (
	"candidate/tools/setting"
	"implement/pay/rpc/orderRPCImp"
	"interface/pay/orderInterfaces"
	"strconv"
	"time"
)

type candidateService struct {
	orderOp orderInterfaces.Operator
}

func NewCandidateService() (*candidateService, error) {
	cs := &candidateService{}
	var err error
	cs.orderOp, err = orderRPCImp.NewOrderRPCImp()
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// CacheCandidate 创建候补订单存入缓存,后将返回点单号给前端, 前端根据订单号支付
func (c candidateService) CacheCandidate(userID, trainId int, date string, passengers []string) (string, error) {
	money := 100
	// 创建订单, 获得外部id
	orderOutsideID, err := c.orderOp.Create(&orderInterfaces.CreateInfo{
		UserID:         int64(userID),
		Money:          strconv.Itoa(money),
		AffairID:       "CAN" + time.Now().Format("2006-01-02-15:04:05.000-") + strconv.Itoa(userID),
		ExpireDuration: 1800,
		CreatedBy:      setting.Server.Name,
	})
	if err != nil {
		return "", err
	}

	// TODO 创建候补存入缓存, 支付完成后, 支付服务通过队列通知候补服务, 候补将该候补信息写入mysql
	return orderOutsideID, err
}
