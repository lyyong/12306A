// @Author LiuYong
// @Created at 2021-02-23
// @Modified at 2021-02-23
package cache

import "fmt"

type OrderCache struct {
	UserID    uint
	OutsideID string
	AffairID  string
}

func (c OrderCache) GetOrderKey() string {
	return fmt.Sprintf("order-%s", c.OutsideID)
}

func (c OrderCache) GetOrdersKey() string {
	return fmt.Sprintf("orders-%d", c.UserID)
}

func (c OrderCache) GetNoFinishOrderKey() string {
	return fmt.Sprintf("order-no-finish-%d", c.UserID)
}
