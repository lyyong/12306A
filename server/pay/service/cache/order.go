// @Author LiuYong
// @Created at 2021-02-23
// @Modified at 2021-02-23
package cache

import "fmt"

type OrderCache struct {
	UserID    uint
	OutsideID string
}

// GetOrderKey 获得一个订单的key 需要OutsideID
func (c OrderCache) GetOrderKey() string {
	return fmt.Sprintf("order-%s", c.OutsideID)
}

// GetOrdersKey 获得一个订单的key 需要UserID
func (c OrderCache) GetOrdersKey() string {
	return fmt.Sprintf("orders-%d", c.UserID)
}

// GetNoFinishOrderKey 获得一个未完成订单的key 需要UserID
func (c OrderCache) GetNoFinishOrderKey() string {
	return fmt.Sprintf("order-no-finish-%d", c.UserID)
}
