/**
 * @Author fzh
 * @Date 2021/2/22
 */
package httpapi

import (
	"common/middleware/token/usertoken"
	"github.com/gin-gonic/gin"
	"net/http"
	"user/service"
	"user/util/resp"
)

type InsertPassengerRequest struct {
	Name              string `json:"name"`
	CertificateType   int    `json:"certificate_type"`
	CertificateNumber string `json:"certificate_number" binding:"certificateNumber"`
	PhoneNumber       string `json:"phone_number" binding:"phoneNumber"`
	PassengerType     int    `json:"passenger_type"`
}

// InsertPassenger godoc
// @Summary 添加乘车人
// @Description 添加乘车人，参数为乘车人信息
// @ID passenger-insert
// @Accept json
// @Produce json
// @Param form body InsertPassengerRequest true "乘车人信息"
// @Success 200 {object} resp.Response
// @Router /passenger [post]
func InsertPassenger(c *gin.Context) {
	req := new(InsertPassengerRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, resp.R(struct{}{}).SetMsg("JSON格式错误"))
		return
	}

	// 从Token中解析出用户ID
	info, ok := usertoken.GetUserInfo(c)
	if !ok {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户未登录"))
		return
	}

	param := &service.PassengerParam{
		UserId:            info.UserId,
		Name:              req.Name,
		CertificateType:   req.CertificateType,
		CertificateNumber: req.CertificateNumber,
		PhoneNumber:       req.PhoneNumber,
		PassengerType:     req.PassengerType,
	}

	if err := service.InsertPassenger(param); err != nil {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人添加失败"))
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人添加成功").SetCode(200))
}

type UpdatePassengerRequest struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	CertificateType   int    `json:"certificate_type"`
	CertificateNumber string `json:"certificate_number" binding:"certificateNumber"`
	PhoneNumber       string `json:"phone_number" binding:"phoneNumber"`
	PassengerType     int    `json:"passenger_type"`
}

// UpdatePassenger godoc
// @Summary 修改乘车人
// @Description 修改乘车人，参数为乘车人信息
// @ID passenger-update
// @Accept json
// @Produce json
// @Param form body UpdatePassengerRequest true "乘车人信息"
// @Success 200 {object} resp.Response
// @Router /passenger [put]
func UpdatePassenger(c *gin.Context) {
	req := new(UpdatePassengerRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, resp.R(struct{}{}).SetMsg("JSON格式错误"))
		return
	}

	// 从Token中解析出用户ID
	info, ok := usertoken.GetUserInfo(c)
	if !ok {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户未登录"))
		return
	}

	param := &service.PassengerParam{
		UserId:            info.UserId,
		PassengerId:       req.ID,
		Name:              req.Name,
		CertificateType:   req.CertificateType,
		CertificateNumber: req.CertificateNumber,
		PhoneNumber:       req.PhoneNumber,
		PassengerType:     req.PassengerType,
	}

	if err := service.UpdatePassenger(param); err != nil {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人修改失败"))
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人修改成功").SetCode(200))
}

type DeletePassengerRequest struct {
	ID uint `json:"id"`
}

// DeletePassenger godoc
// @Summary 删除乘车人
// @Description 删除乘车人，参数为乘车人信息
// @ID passenger-delete
// @Accept json
// @Produce json
// @Param form body DeletePassengerRequest true "乘车人ID"
// @Success 200 {object} resp.Response
// @Router /passenger [delete]
func DeletePassenger(c *gin.Context) {
	req := new(DeletePassengerRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, resp.R(struct{}{}).SetMsg("JSON格式错误"))
		return
	}

	// 从Token中解析出用户ID
	info, ok := usertoken.GetUserInfo(c)
	if !ok {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户未登录"))
		return
	}

	param := &service.PassengerParam{
		UserId:      info.UserId,
		PassengerId: req.ID,
	}

	if err := service.DeletePassenger(param); err != nil {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人删除失败"))
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人删除成功").SetCode(200))
}

type ListPassengerResponse struct {
	Passenger []*ListPassengerResponseData `json:"passenger"`
}

type ListPassengerResponseData struct {
	Id                uint   `json:"id"`
	Name              string `json:"name"`
	CertificateType   int    `json:"certificateType"`
	CertificateNumber string `json:"certificateNumber"`
	PhoneNumber       string `json:"phoneNumber"`
	PassengerType     int    `json:"passengerType"`
}

// ListPassenger godoc
// @Summary 查询乘车人
// @Description 查询乘车人
// @ID passenger-list
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response
// @Router /passenger [get]
func ListPassenger(c *gin.Context) {
	// 从Token中解析出用户ID
	info, ok := usertoken.GetUserInfo(c)
	if !ok {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户未登录"))
		return
	}

	response := new(ListPassengerResponse)
	if passengers, err := service.ListPassenger(info.UserId); err != nil {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人修改失败"))
		return
	} else {
		for _, p := range passengers {
			data := &ListPassengerResponseData{
				Id:                p.Id,
				Name:              p.Name,
				CertificateType:   p.CertificateType,
				CertificateNumber: p.CertificateNumber,
				PhoneNumber:       p.PhoneNumber,
				PassengerType:     p.PassengerType,
			}
			response.Passenger = append(response.Passenger, data)
		}
	}

	c.JSON(http.StatusOK, resp.R(response).SetCode(200))
}
