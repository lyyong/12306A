// @Author: KongLingWen
// @Created at 2021/2/22
// @Modified at 2021/2/22

package client

import "rpc/ticket/proto/ticketRPC"

type TicketRPCClient struct {
	pbClient ticketRPC.TicketServiceClient
}

var client *TicketRPCClient

