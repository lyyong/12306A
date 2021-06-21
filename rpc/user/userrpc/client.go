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

// Deprecated: Use NewClientWithTarget instead.
func NewClient() *Client {
	return NewClientWithTarget(targetService)
}

// 创建RPC客户端
func NewClientWithTarget(target string) *Client {
	conn, err := rpc_manage.NewGRPCClientConn(target)
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

type Passenger struct {
	Id                uint
	Name              string
	CertificateType   int
	CertificateNumber string
	PassengerType     int
}

// 根据用户ID查询乘车人列表
func (c *Client) ListPassenger(id uint) ([]*Passenger, error) {
	res, err := c.client.ListPassenger(context.Background(), &userpb.ListPassengerRequest{Id: uint32(id)})
	if err != nil {
		return nil, err
	}
	var list []*Passenger
	for _, p := range res.Passenger {
		passenger := &Passenger{
			Id:                uint(p.Id),
			Name:              p.Name,
			CertificateType:   int(p.CertificateType),
			CertificateNumber: p.CertificateNumber,
			PassengerType:     int(p.PassengerType),
		}
		list = append(list, passenger)
	}
	return list, nil
}
