// @Author LiuYong
// @Created at 2021-01-29
// @Modified at 2021-01-29
package server_find

import (
	"common/tools/logging"
	"errors"
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"time"
)

// Client 作为consul服务的封装, 包含了服务发现和注销服务
// server_find.Client 将会是一个全局单例.
// 首先是由 server_find.Register() 初始化 server_find.Client.
// 然后可以直接通过 server_find.GetClient()获取.
// 使用 consul.HaveClient() 判断是否进行了 server_find.Client 的初始化
// 服务器关闭后要使用 server_find.DeRegister()关闭.
type Client struct {
	cli       *consulApi.Client
	name      string // 服务名称
	host      string // 服务的地址
	port      int    // 服务的端口号
	serviceID string // 服务的ID
	target    string // consul的地址,例如"localhost:8500"
	interval  int    // consul的地址,例如"localhost:8500"
	ttl       int    // 注册信息的缓存时间, 如果ttl过时前没有得到updateTTL则注册服务信息将被抛弃, 单位秒
}

func (c *Client) Target() string {
	return c.target
}

func (c *Client) Agent() *consulApi.Agent {
	return c.cli.Agent()
}

func (c *Client) Health() *consulApi.Health {
	return c.cli.Health()
}

// newClient 注册服务到consul
// name - 服务名称
// host - 服务的地址
// port - 服务的端口号
// serviceID - 服务的ID
// target - consul的地址,例如"localhost:8500"
// interval - //consul的地址,例如"localhost:8500"
// ttl - 注册信息的缓存时间, 如果ttl过时前没有得到updateTTL则注册服务信息将被抛弃, 单位秒
func newClient(name string, host string, port int, serviceID string, target string, interval int, ttl int) (*Client, error) {
	if interval > ttl {
		return nil, errors.New("interval大于ttl")
	}
	return &Client{
		cli:       nil,
		name:      name,
		host:      host,
		port:      port,
		serviceID: serviceID,
		target:    target,
		interval:  interval,
		ttl:       ttl,
	}, nil
}

// register 注册服务
func (c *Client) register() error {
	if c.cli != nil {
		return errors.New("重复注册服务")
	}
	// consul客户端配置
	conf := &consulApi.Config{Scheme: "http", Address: c.target}
	// 创建consul 连接客户端
	var err error
	c.cli, err = consulApi.NewClient(conf)
	if err != nil {
		return fmt.Errorf("创建consul client客户端错误: %v", err)
	}

	// 更新ttl
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(c.interval))
		for {
			<-ticker.C
			err := c.cli.Agent().UpdateTTL(c.serviceID, "", consulApi.HealthPassing)
			if err != nil {
				logging.Error("更新服务的ttl错误: %v", err)
				break
			}
		}
	}()

	// 注册服务信息
	err = c.cli.Agent().ServiceRegister(&consulApi.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    c.name,
		Port:    c.port,
		Address: c.host,
	})
	if err != nil {
		return fmt.Errorf("注册服务%s 出错 %v", c.name, err)
	}

	// 检查是否注册完成
	check := consulApi.AgentServiceCheck{
		TTL:    fmt.Sprintf("%ds", c.ttl),
		Status: consulApi.HealthPassing,
	}
	err = c.cli.Agent().CheckRegister(&consulApi.AgentCheckRegistration{
		ID:                c.serviceID,
		Name:              c.name,
		ServiceID:         c.serviceID,
		AgentServiceCheck: check,
	})
	if err != nil {
		return fmt.Errorf("注册服务检查出错: %v", err)
	}
	return nil
}

// deRegister 注销服务
func (c *Client) deRegister() error {
	if c.cli == nil {
		return errors.New("服务未注册")
	}
	defer func() {
		c.cli = nil
	}()
	err := c.cli.Agent().ServiceDeregister(c.serviceID)
	if err != nil {
		return fmt.Errorf("注销服务出错%v", err)
	}
	err = c.cli.Agent().CheckDeregister(c.serviceID)
	if err != nil {
		return fmt.Errorf("检查注销服务出错: %v", err)
	}
	return nil
}

// isRegister 是否注册了服务
func (c *Client) isRegister() bool {
	return c.cli != nil
}
