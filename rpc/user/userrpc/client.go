/**
 * @Author fzh
 * @Date 2021/2/21
 */
package userrpc

import (
	"common/rpc_manage"
	"common/tools/logging"
	"context"
	"rpc/user/userpb"
)

const targetService = "localhost:8224"

type Client struct {
	client userpb.UserServiceClient
}

func NewClient() *Client {
	conn, err := rpc_manage.NewGRPCClientConn(targetService)
	if err != nil {
		logging.Fatal("用户RPC客户端创建失败")
	}
	rpcClient := userpb.NewUserServiceClient(conn)
	return &Client{client: rpcClient}
}

type User struct {
	Id                uint
	Username          string
	State             int
	CertificateType   int
	Name              string
	CertificateNumber string
	PhoneNumber       string
	Email             string
	PassengerType     int
}

func (c *Client) GetUser(id uint) (*User, error) {
	res, err := c.client.GetUser(context.Background(), &userpb.UserRequest{Id: uint32(id)})
	if err != nil {
		return nil, err
	}
	user := &User{
		Id:                uint(res.Id),
		Username:          res.Username,
		State:             int(res.State),
		CertificateType:   int(res.CertificateType),
		Name:              res.Name,
		CertificateNumber: res.CertificateNumber,
		PhoneNumber:       res.PhoneNumber,
		Email:             res.Email,
		PassengerType:     int(res.PassengerType),
	}
	return user, nil
}
