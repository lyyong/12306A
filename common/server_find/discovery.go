// @Author LiuYong
// @Created at 2021-02-01
// @Modified at 2021-02-01
package server_find

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/resolver"
	"strconv"
)

// ServiceDiscoveryAboutBalance 负载均衡相关的服务查找
// 主要是查询服务注册服务器查找自己要找的服务
type ServiceDiscoveryAboutBalance struct {
	consulClient      *Client
	targetServiceName string
}

// NewServiceDiscoveryAboutBalance 通过consul.Client创建一个负载均衡查询功能
// Client必须是已经Register过的,
// tgSN - targetServiceName 需要进行rpc请求的服务器名称
func NewServiceDiscoveryAboutBalance(cClient *Client, tgSN string) (*ServiceDiscoveryAboutBalance, error) {
	if !cClient.isRegister() {
		return nil, errors.New("该consulClient没有进行服务注册")
	}
	return &ServiceDiscoveryAboutBalance{consulClient: cClient, targetServiceName: tgSN}, nil
}

// Build 功能主要是通过target创建一个resolver, 会在grpc.Dial()时执行
// 系统默认的resolver有三种:
// dns.dnsResolver 通过域名解析获得服务地址;
// manual.Resolver 手动设置服务地址;
// passthroughResolver Dial()中第一个参数直接作为服务地址;
//
// target - 让我们知道scheme和endpoint等信息,让我们可以判断不同的请求用不同的服务发现管理
// 该函数需要把resolver.ClientConn的状态进行更新, 主要就是更新地址,地址为一个切片,
// 更新cc的地址后就可以有选择
func (s ServiceDiscoveryAboutBalance) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 向consul请求获取健康的服务名为targetServiceName的服务器信息
	serviceEntrys, _, err := s.consulClient.Health().Service(s.targetServiceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	addrs := make([]resolver.Address, len(serviceEntrys))
	// 收集地址
	for _, se := range serviceEntrys {
		addrs = append(addrs, resolver.Address{
			Addr: se.Service.Address + ":" + strconv.Itoa(se.Service.Port),
		})
	}
	// 更新地址
	cc.UpdateState(resolver.State{Addresses: addrs})
	return s, nil
}

func (s ServiceDiscoveryAboutBalance) Scheme() string {
	return SCHEME
}

func (s ServiceDiscoveryAboutBalance) ResolveNow(options resolver.ResolveNowOptions) {
	fmt.Println("ResolveNew")
}

func (s ServiceDiscoveryAboutBalance) Close() {
	fmt.Println("Close")
}
