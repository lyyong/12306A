// @Author LiuYong
// @Created at 2021-01-29
// @Modified at 2021-01-29
package consul

import (
	"errors"
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"pay/tools/logging"
	"time"
)

var client *consulApi.Client
var sID string // 服务ID

// Register 注册服务到consul
// name - 服务名称
// host - 服务的地址
// port - 服务的端口号
// target - consul的地址,例如"localhost:8500"
// interval - 刷新ttl的时间间隔, 单位秒
// ttl - 注册信息的缓存时间, 如果ttl过时前没有得到updateTTL则注册服务信息将被抛弃, 单位秒
func Register(name string, host string, port int, serviceID string, target string, interval int, ttl int) error {
	if interval > ttl {
		return errors.New("interval大于ttl")
	}
	if client != nil {
		return errors.New("重复注册服务")
	}
	// consul客户端配置
	conf := &consulApi.Config{Scheme: "http", Address: target}
	// 创建consul 连接客户端
	var err error
	client, err = consulApi.NewClient(conf)
	if err != nil {
		return fmt.Errorf("创建consul client客户端错误: %v", err)
	}

	// 服务ID, 通过名称-地址-端口号
	sID = serviceID

	// 更新ttl
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(interval))
		for {
			<-ticker.C
			err := client.Agent().UpdateTTL(sID, "", consulApi.HealthPassing)
			if err != nil {
				logging.Error("更新服务的ttl错误", err)
			}
		}
	}()

	// 注册服务信息
	err = client.Agent().ServiceRegister(&consulApi.AgentServiceRegistration{
		ID:      sID,
		Name:    name,
		Port:    port,
		Address: host,
	})
	if err != nil {
		return fmt.Errorf("注册服务%s 出错 %v", name, err)
	}

	// 检查是否注册完成
	check := consulApi.AgentServiceCheck{
		TTL:    fmt.Sprintf("%ds", ttl),
		Status: consulApi.HealthPassing,
	}
	err = client.Agent().CheckRegister(&consulApi.AgentCheckRegistration{
		ID:                sID,
		Name:              name,
		ServiceID:         sID,
		AgentServiceCheck: check,
	})
	if err != nil {
		return fmt.Errorf("注册服务检查出错: %v", err)
	}
	return nil
}

func DeRegister() error {
	if client == nil {
		return errors.New("服务未注册")
	}
	err := client.Agent().ServiceDeregister(sID)
	if err != nil {
		return fmt.Errorf("注销服务出错%v", err)
	}
	err = client.Agent().CheckDeregister(sID)
	if err != nil {
		return fmt.Errorf("检查注销服务出错: %v", err)
	}
	return nil
}
