package service

import (
	"common/tools/logging"
	orderClient "rpc/pay/client/orderRPCClient"
	ticketPoolClient "rpc/ticketPool/Client"
	"ticket/utils/setting"
)

var ts *TicketService

type TicketService struct {
	orderCli	*orderClient.OrderRPCClient
	tpCli		*ticketPoolClient.TPRPCClient
}

func init(){
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
	ts.orderCli, err = orderClient.NewClientWithMQHost(setting.Kafka.Host)
	ts.orderCli.SetDealPayOK(payOK)
	if err != nil {
		return nil, err
	}
	ts.tpCli,err = ticketPoolClient.NewClient()
	if err != nil {
		return nil, err
	}
	return ts, nil
}







