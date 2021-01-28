// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package orderInterfaces

// CreateInfo 创建订单需要提供的信息  创建的订单都是未支付的
type CreateInfo struct {
	UserID         int64  `json:"user_id"`          // 创建的用户id
	Money          string `json:"money"`            // 支付的金额
	AffairID       string `json:"affair_id"`        // 事务id 可能是redis存储与该订单相关信息的key
	ExpireDuration int32  `json:"expire_duration"`  // 过期时间 单位秒
	OrderOutsideID string `json:"order_outside_id"` // 外部id 暂时可以为空或者一个随机数
}

type Creator interface {
	// 创建订单
	Create(info *CreateInfo) error
}
