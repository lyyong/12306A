package service

import (
	"common/tools/logging"
	orderClient "rpc/pay/client/orderRPCClient"
	ticketPoolClient "rpc/ticketPool/Client"
	userClient "rpc/user/userrpc"
	"ticket/utils/setting"
)

var ts *TicketService

type TicketService struct {
	orderCli *orderClient.OrderRPCClient
	tpCli    *ticketPoolClient.TPRPCClient
	userCli  *userClient.Client
}

func init() {
	logging.Info("Init Ticket Service")
	var err error
	ts, err = NewTicketService()
	if err != nil {
		logging.Fatal("Fail to init Service:", err)
	}
}

func NewTicketService() (*TicketService, error) {
	ts := &TicketService{}
	var err error
	ts.orderCli, err = orderClient.NewClientWithTargetAndMQHost(setting.RpcTarget.Order, setting.Kafka.Host)
	ts.orderCli.SetDealPayOK(payOK)
	if err != nil {
		return nil, err
	}
	ts.tpCli, err = ticketPoolClient.NewClientWithTarget(setting.RpcTarget.TicketPool)
	if err != nil {
		return nil, err
	}
	ts.userCli = userClient.NewClientWithTarget(setting.RpcTarget.User)
	return ts, nil
}
