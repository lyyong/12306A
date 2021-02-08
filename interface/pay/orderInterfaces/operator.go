// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package orderInterfaces

type Operator interface {
	// 创建订单
	Create(info *CreateInfo) (string, error)
	Read(userID int64) ([]*Info, error)
	UpdateState(orderOutsideID string, state int32) error
	UpdateStateWithRelativeOrder(orderOutsideID string, state int32, relativeOutsideID string) error
}
