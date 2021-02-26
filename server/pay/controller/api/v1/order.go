// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package v1

import "github.com/gin-gonic/gin"

type UserFinishOrdersSend struct {
}

// @Summary 用户获取自己的订单信息
// @Description
// @Accept json
// @Produce json
// @Param token header string true "认证信息"
// @Param wantPayR body v1.payOKAbbRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /ok/abb [get]
func GetUserFinishOrders(c *gin.Context) {

}
