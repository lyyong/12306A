// @Author LiuYong
// @Created at 2021-02-01
// @Modified at 2021-02-01
package server_find

import (
	"common/tools/logging"
	"errors"
	"strconv"
)

var consulClient *Client

// Register 注册服务
// name - 服务名称
// host - 服务的地址
// port - 服务的端口号
// serviceID - 服务的ID
// target - consul的地址,例如"localhost:8500"
// interval - 心跳间隙
// ttl - 注册信息的缓存时间, 如果ttl过时前没有得到updateTTL则注册服务信息将被抛弃, 单位秒
func Register(name, host, port, serviceID, target string, interval, ttl int) error {
	if consulClient != nil {
		return errors.New("重复注册服务")
	}
	var err error
	pt, err := strconv.Atoi(port)
	if err != nil {
		consulClient = nil
		return err
	}
	consulClient, err = newClient(name, host, pt, serviceID, target, interval, ttl)
	if err != nil {
		consulClient = nil
		return err
	}
	logging.Info("开启服务发现功能, 服务ID: ", serviceID, " consul-host: ", target)
	err = consulClient.register()
	if err != nil {
		consulClient = nil
		logging.Error("服务发现功能开启失败")
		return err
	}
	logging.Info("服务注册完成")
	return nil
}

// DeRegister 注销服务
func DeRegister() error {
	if consulClient == nil {
		return errors.New("未注册服务, 注销失败")
	}
	if err := consulClient.deRegister(); err != nil {
		return err
	}
	return nil
}

// IsRegister 检查服务是否已经注册
func IsRegister() bool {
	return consulClient != nil
}

// GetClient 获取 server_find.Client 单例
func GetClient() (*Client, error) {
	if consulClient == nil {
		return nil, errors.New("未进行服务注册")
	}
	return consulClient, nil
}
