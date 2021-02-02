// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCImp

import (
	"common/router_tracer"
	"common/server_find"
	"interface/pay/orderInterfaces"
	"testing"
)

func TestOrderRPCImp_Create(t *testing.T) {
	err := server_find.Register("TestOrderRPCImp", "localhost", "9001", "TestOrderRPCImp-localhost-9001", "localhost:8500", 15, 20)
	if err != nil {
		t.Error(err)
	}
	err = router_tracer.SetupByHttp("TestOrderRPCImp", "localhost", "9001", "http://localhost:9411/api/v2/spans")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		server_find.DeRegister()
		router_tracer.Close()
	}()
	rpcImp, err := NewOrderRPCImp()
	if err != nil {
		t.Error(err)
	}
	for i := 100; i > 0; i-- {
		ret := rpcImp.Create(&orderInterfaces.CreateInfo{
			UserID:         1,
			Money:          "123",
			AffairID:       "asd",
			ExpireDuration: 0,
			OrderOutsideID: "123zxcz",
		})
		t.Log(ret)
	}
}
