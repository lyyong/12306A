// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCClient

import (
	"rpc/pay/proto/orderRPCpb"
	"testing"
)

func TestCreate(t *testing.T) {
	client, err := NewClientWithHttpTracer("test", "localhost", "9000", "http://localhost:9411/api/v2/spans")
	defer client.Close()
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Create(&orderRPCpb.CreateInfo{
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
}

func TestRead(t *testing.T) {
	client, err := NewClient()
	defer client.Close()
	resp, err := client.Read(&orderRPCpb.SearchInfo{UserID: 1})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
