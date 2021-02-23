// @Author LiuYong
// @Created at 2021-01-30
// @Modified at 2021-01-30
package router_tracer

import (
	"fmt"
	"github.com/openzipkin/zipkin-go"
	zgrp "github.com/openzipkin/zipkin-go/reporter"
)

// router_tracer.Client 将会是一个全局单例.
// 首先是由zipkin.SetupByHttp() 初始化zipkin.Client.
// 然后可以直接通过zipkin.GetClient()获取.
// 使用zipkin.IsTracing() 判断是否进行了zipkin的初始化
// 服务器关闭后要使用zipkin.Close()关闭.
type Client struct {
	httpReporter  *zgrp.Reporter
	kafkaReporter *zgrp.Reporter
	serviceName   string
	serviceHost   string
	servicePort   string
	zkHttp        string
	tracer        *zipkin.Tracer
}

func (c *Client) Tracer() *zipkin.Tracer {
	return c.tracer
}

// newClient 创建一个zipkin客户端, 用来连接zipkin服务器, 给服务器汇报信息
// serviceName - 当前服务的名称
// serviceHost - 当前服务的host 不包括端口
// servicePort - 当前服务的端口
func newClient(serviceName string, serviceHost string, servicePort string) *Client {
	return &Client{
		httpReporter:  nil,
		kafkaReporter: nil,
		serviceName:   serviceName,
		serviceHost:   serviceHost,
		servicePort:   servicePort,
		zkHttp:        "",
		tracer:        nil,
	}
}

func (c *Client) setup(zkReporter *zgrp.Reporter) error {
	// 创建一个服务节点, 作为链路上的一个节点
	endpoint, err := zipkin.NewEndpoint(c.serviceName, c.serviceHost+":"+c.servicePort)
	if err != nil {
		return fmt.Errorf("创建节点出错: %v", err)
	}

	tracer, err := zipkin.NewTracer(*zkReporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("创建追踪器错误: %v", err)
	}

	// 把追踪器包装成opentracing
	// t := zipkinot.Wrap(nativeTracer)
	c.tracer = tracer
	return nil
}

func (c *Client) closeHttp() {
	if c.httpReporter != nil {
		(*c.httpReporter).Close()
	}
}

func (c *Client) closeKafka() {
	if c.kafkaReporter != nil {
		(*c.kafkaReporter).Close()
	}
}
