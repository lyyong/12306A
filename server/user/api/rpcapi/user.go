/**
 * @Author fzh
 * @Date 2021/2/21
 */
package rpcapi

import (
	"context"
	"rpc/user/userpb"
	"user/service"
)

type UserService struct {
}

func (s *UserService) GetUser(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	id := uint(req.GetId())
	user, err := service.GetUser(id)
	if err != nil {
		return nil, err
	}
	response := &userpb.UserResponse{
		Id:                uint32(user.ID),
		Username:          user.Username,
		State:             int32(user.State),
		CertificateType:   int32(user.CertificateType),
		Name:              user.Name,
		CertificateNumber: user.CertificateNumber,
		PhoneNumber:       user.PhoneNumber,
		Email:             user.Email,
		PassengerType:     int32(user.PassengerType),
	}
	return response, nil
}
