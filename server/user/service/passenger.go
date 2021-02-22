/**
 * @Author fzh
 * @Date 2021/2/22
 */
package service

import (
	"common/tools/logging"
	. "user/global/database"
	"user/model"
)

type PassengerParam struct {
	UserId uint
	Data   []*PassengerParamData
}

type PassengerParamData struct {
	PassengerId       uint
	Name              string
	CertificateType   int
	CertificateNumber string
	PassengerType     int
}

// 添加乘车人
func InsertPassenger(param *PassengerParam) error {
	var passengers []*model.Passenger
	for _, data := range param.Data {
		passenger := &model.Passenger{
			Name:              data.Name,
			CertificateType:   data.CertificateType,
			CertificateNumber: data.CertificateNumber,
			PassengerType:     data.PassengerType,
			State:             0,
		}
		passengers = append(passengers, passenger)
	}
	err := model.InsertPassenger(DB, param.UserId, passengers)
	if err != nil {
		logging.Error("添加乘车人出错")
		// TODO: 返回错误
	}
	return nil
}

// 更新乘车人
func UpdatePassenger(param *PassengerParam) error {
	var passengers []*model.Passenger
	for _, data := range param.Data {
		passenger := &model.Passenger{
			Name:              data.Name,
			CertificateType:   data.CertificateType,
			CertificateNumber: data.CertificateNumber,
			PassengerType:     data.PassengerType,
			State:             0,
		}
		passengers = append(passengers, passenger)
	}
	err := model.UpdatePassenger(DB, param.UserId, passengers)
	if err != nil {
		logging.Error("修改乘车人出错")
		// TODO: 返回错误
	}
	return nil
}

type PassengerRecord struct {
	Id                uint
	Name              string
	CertificateType   int
	CertificateNumber string
	PassengerType     int
	State             int
}

func ListPassenger(userId uint) ([]*PassengerRecord, error) {
	passengers, err := model.ListPassenger(DB, userId)
	if err != nil {
		logging.Error("查询乘车人出错")
		// TODO: 返回错误
	}
	var result []*PassengerRecord
	for _, p := range passengers {
		record := &PassengerRecord{
			Id:                p.ID,
			Name:              p.Name,
			CertificateType:   p.CertificateType,
			CertificateNumber: p.CertificateNumber,
			PassengerType:     p.PassengerType,
			State:             p.State,
		}
		result = append(result, record)
	}
	return result, nil
}
