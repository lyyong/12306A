// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package orderMockImp

import (
	"errors"
	"interface/pay/orderInterfaces"
)

type MockOrder struct{}

func (o *MockOrder) UpdateState(orderOutsideID string, state int32) error {
	panic("implement me")
}

func (o *MockOrder) UpdateStateWithRelativeOrder(orderOutsideID string, state int32, relativeOutsideID string) error {
	panic("implement me")
}

func (*MockOrder) Create(info *orderInterfaces.CreateInfo) (string, error) {
	if info.UserID < 0 {
		return "", errors.New("不合法的userID")
	}
	return "", nil
}

func (*MockOrder) Read(userID int64) (*orderInterfaces.Info, error) {
	if userID < 0 {
		return nil, errors.New("不合法的userID")
	}
	return &orderInterfaces.Info{}, nil
}
