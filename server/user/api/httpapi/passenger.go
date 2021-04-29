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
	Data []InsertPassengerRequestData
}

type InsertPassengerRequestData struct {
	Name              string
	CertificateType   int
	CertificateNumber string
	PassengerType     int
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

	param := new(service.PassengerParam)
	param.UserId = info.UserId
	for _, p := range req.Data {
		data := &service.PassengerParamData{
			Name:              p.Name,
			CertificateType:   p.CertificateType,
			CertificateNumber: p.CertificateNumber,
			PassengerType:     p.PassengerType,
		}
		param.Data = append(param.Data, data)
	}

	if err := service.InsertPassenger(param); err != nil {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人添加失败"))
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人添加成功").SetCode(200))
}

type UpdatePassengerRequest struct {
	Data []UpdatePassengerRequestData
}

type UpdatePassengerRequestData struct {
	Name              string
	CertificateType   int
	CertificateNumber string
	PassengerType     int
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

	param := new(service.PassengerParam)
	param.UserId = info.UserId
	for _, p := range req.Data {
		data := &service.PassengerParamData{
			Name:              p.Name,
			CertificateType:   p.CertificateType,
			CertificateNumber: p.CertificateNumber,
			PassengerType:     p.PassengerType,
		}
		param.Data = append(param.Data, data)
	}

	if err := service.UpdatePassenger(param); err != nil {
		c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人修改失败"))
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("乘车人修改成功").SetCode(200))
}

type ListPassengerResponse struct {
	Passenger []*ListPassengerResponseData `json:"passenger"`
}

type ListPassengerResponseData struct {
	Id                uint   `json:"id"`
	Name              string `json:"name"`
	CertificateType   int    `json:"certificateType"`
	CertificateNumber string `json:"certificateNumber"`
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
				PassengerType:     p.PassengerType,
			}
			response.Passenger = append(response.Passenger, data)
		}
	}

	c.JSON(http.StatusOK, resp.R(response).SetCode(200))
}
