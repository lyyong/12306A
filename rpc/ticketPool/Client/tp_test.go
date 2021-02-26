// @Author LiuYong
// @Created at 2021-02-26
// @Modified at 2021-02-26
package Client

import (
	"rpc/ticketPool/proto/ticketPoolRPC"
	"testing"
)

func Test_GetTicketNumber(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	res, err := cli.GetTicketNumber(&ticketPoolRPC.GetTicketNumberRequest{
		TrainId:        []uint32{2},
		StartStationId: 1,
		DestStationId:  4,
		Date:           "2021-02-28",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(res.TrainsTicketInfo[0])
}
