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
		Date: "2021-02-28",
		Condition: []*ticketPoolRPC.GetTicketNumberRequest_Condition{
			&ticketPoolRPC.GetTicketNumberRequest_Condition{
				TrainId:        17416,
				StartStationId: 3337,
				DestStationId:  3326,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(res.TrainsTicketInfo[0])
}

func Test_BuyTicket(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	res, err := cli.GetTicket(&ticketPoolRPC.GetTicketRequest{
		TrainId:        17416,
		StartStationId: 3337,
		DestStationId:  3326,
		Date:           "2021-02-28",
		Passengers: []*ticketPoolRPC.PassengerInfo{
			&ticketPoolRPC.PassengerInfo{
				PassengerId:   1,
				PassengerName: "123",
				SeatTypeId:    0,
				ChooseSeat:    "A",
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(res.Tickets[0])
}
