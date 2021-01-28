// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCClient

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"rpc/pay/proto/orderRPCpb"
)

var client *orderRPCpb.OrderRPCServiceClient
var conn *grpc.ClientConn

const server_port = "8082"

// InitClient 初始化客户端连接
func InitClient() error {
	var err error
	if client != nil {
		return errors.New("重复创建OrderRPCServiceClient")
	}
	if conn == nil {
		conn, err = grpc.Dial(":"+server_port, grpc.WithInsecure())
		if err != nil {
			return err
		}
	}
	tclient := orderRPCpb.NewOrderRPCServiceClient(conn)
	client = &tclient
	return nil
}

func CloseClient() error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}

func Create(info *orderRPCpb.CreateInfo) (*orderRPCpb.Error, error) {
	if client == nil {
		InitClient()
	}
	tclient := *client
	resp, err := tclient.Create(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Read(info *orderRPCpb.SearchInfo) (*orderRPCpb.Info, error) {
	if client == nil {
		InitClient()
	}
	tclient := *client
	resp, err := tclient.Read(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}
