// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package orderInterfaces

type Operator interface {
	// 创建订单
	Create(info *CreateInfo) error
	Read(userID int64) (*Info, error)
}
