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
	UserId            uint
	PassengerId       uint
	Name              string
	CertificateType   int
	CertificateNumber string
	PhoneNumber       string
	PassengerType     int
}

// 添加乘车人
func InsertPassenger(param *PassengerParam) error {
	passenger := &model.Passenger{
		Name:              param.Name,
		CertificateType:   param.CertificateType,
		CertificateNumber: param.CertificateNumber,
		PhoneNumber:       param.PhoneNumber,
		PassengerType:     param.PassengerType,
		State:             0,
	}

	err := model.InsertPassenger(DB, param.UserId, passenger)
	if err != nil {
		logging.Error("添加乘车人出错")
		// TODO: 返回错误
	}
	return nil
}

// 更新乘车人
func UpdatePassenger(param *PassengerParam) error {
	passenger := &model.Passenger{
		Name:              param.Name,
		CertificateType:   param.CertificateType,
		CertificateNumber: param.CertificateNumber,
		PhoneNumber:       param.PhoneNumber,
		PassengerType:     param.PassengerType,
		State:             0,
	}
	passenger.ID = param.PassengerId

	err := model.UpdatePassenger(DB, param.UserId, passenger)
	if err != nil {
		logging.Error("修改乘车人出错")
		// TODO: 返回错误
	}
	return nil
}

func DeletePassenger(param *PassengerParam) error {
	passenger := &model.Passenger{}
	passenger.ID = param.PassengerId

	err := model.DeletePassenger(DB, param.UserId, passenger)
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
	PhoneNumber       string
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
			PhoneNumber:       p.PhoneNumber,
			PassengerType:     p.PassengerType,
			State:             p.State,
		}
		result = append(result, record)
	}
	return result, nil
}
