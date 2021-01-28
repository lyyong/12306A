// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package orderMockImp

import (
	"interface/pay/orderInterfaces"
	"testing"
)

func TestMockOrder_Create(t *testing.T) {
	var op orderInterfaces.Operator
	op = &MockOrder{}
	err := op.Create(&orderInterfaces.CreateInfo{UserID: 1})
	if err != nil {
		t.Error(err)
	}

}
