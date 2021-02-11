// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderInterfaces

type Updator interface {
	// 更新订单的状态
	UpdateState(orderOutsideID string, state int32) error
	// 更新订单的状态并且添加相关的订单
	UpdateStateWithRelativeOrder(orderOutsideID string, state int32, relativeOutsideID string) error
}
