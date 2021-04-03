// @Author: KongLingWen
// @Created at 2021/4/4
// @Modified at 2021/4/4

package rpc

import (
	"context"
	"fmt"
	pb "rpc/ticket/proto/ticketRPC"
	"testing"
)

func TestTicketServer_GetTicketByOrdersId(t *testing.T) {
	ordersId := make([]string,1)
	ordersId[0] = ""
	ts := &TicketServer{}
	res, err := ts.GetTicketByOrdersId(context.Background(),&pb.GetTicketByOrdersIdRequest{OrdersId: ordersId})
	if err != nil {
		t.Error()
	}
	fmt.Println(res)
}
