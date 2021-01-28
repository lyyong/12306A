// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCImp

import (
	"interface/pay/orderInterfaces"
	"testing"
)

func TestOrderRPCImp_Create(t *testing.T) {
	rpcImp, err := NewOrderRPCImp()
	if err != nil {
		t.Error(err)
	}
	ret := rpcImp.Create(&orderInterfaces.CreateInfo{
		UserID:         1,
		Money:          "123",
		AffairID:       "asd",
		ExpireDuration: 0,
		OrderOutsideID: "123zxcz",
	})
	t.Log(ret)
	rpcImp.Close()
}
