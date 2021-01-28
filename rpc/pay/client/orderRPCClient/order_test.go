// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCClient

import (
	"rpc/pay/proto/orderRPCpb"
	"testing"
)

func TestCreate(t *testing.T) {
	InitClient()
	resp, err := Create(&orderRPCpb.CreateInfo{
		UserID:         1,
		Money:          "20",
		AffairID:       "123asd",
		ExpireDuration: 0,
		OrderOutsideID: "123zxdc",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.Content)
	defer CloseClient()
}

func TestRead(t *testing.T) {
	InitClient()
	resp, err := Read(&orderRPCpb.SearchInfo{UserID: 1})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	defer CloseClient()
}
