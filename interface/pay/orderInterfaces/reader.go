// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package orderInterfaces

// Info 订单的信息
type Info struct {
	UserID         int64  `json:"user_id"`          // 创建的用户id
	Money          string `json:"money"`            // 支付的金额
	AffairID       string `json:"affair_id"`        // 事务id 可能是redis存储与该订单相关信息的key
	ExpireDuration int32  `json:"expire_duration"`  // 过期时间 单位秒
	OrderOutsideID string `json:"order_outside_id"` // 外部id 暂时可以为空或者一个随机数
	State          int32  `json:"state"`            // 支付状态 0未支付 1已支付
}

type Reader interface {
	// 获取订单信息
	Read(userID int64) ([]*Info, error)
}
