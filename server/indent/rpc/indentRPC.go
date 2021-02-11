package rpc

import (
	"context"
	pb "rpc/indent/proto/indentRPC"
)

type IndentServer struct {

}

func (is *IndentServer) CreateIndent(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	// 生成外部订单号

	// 写入redis，设置过期时间

	// 返回外部订单号

	return &pb.CreateResponse{
		IndentOuterId: "E00001",
	}, nil
}

func (is *IndentServer) PayIndent(ctx context.Context, in *pb.PayRequest) (*pb.PayResponse, error) {
	// 根据用户 id 读 redis ，如果有数据则写入db, 如果已过期，返回 false 支付模块退款
	return &pb.PayResponse{
		IsOk: true,
	},nil
}

func (is *IndentServer) HasUnfinishedIndent(ctx context.Context, in *pb.UnfinishedRequest) (*pb.UnfinishedResponse, error) {
	// 根据用户 id 读 redis 判断是否有未完成的订单

	return &pb.UnfinishedResponse{
		HasUnfinishedIndent: false,
	}, nil

}