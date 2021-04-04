// @Author LiuYong
// @Created at 2021-04-04
// @Modified at 2021-04-04
package client

import (
	"rpc/ticket/proto/ticketRPC"
	"testing"
)

func TestTicketRPCClient_GetTicketByOrdersId(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	res, err := cli.GetTicketByOrdersId(&ticketRPC.GetTicketByOrdersIdRequest{
		OrdersId: []string{"ticket2021-03-04-12:57:14.529-2"},
	})
	t.Log(res)
}
