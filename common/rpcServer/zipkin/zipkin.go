// @Author LiuYong
// @Created at 2021-01-30
// @Modified at 2021-01-30
package zipkin

import (
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zgrp "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

type Client struct {
	httpReporter  *zgrp.Reporter
	kafkaReporter *zgrp.Reporter
	serviceName   string
	serviceHost   string
	servicePort   string
	zkHttp        string
	tracer        *opentracing.Tracer
}

func (c *Client) Tracer() *opentracing.Tracer {
	return c.tracer
}

func NewClient(serviceName string, serviceHost string, servicePort string) *Client {
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

func (c *Client) SetupByHttp(zkHttp string) error {
	// 创建一个报告器, 用来想zipkin汇报信息
	if c.httpReporter != nil {
		return errors.New("重复创建zipkin-http")
	}
	rp := zipkinhttp.NewReporter(zkHttp)
	c.httpReporter = &rp

	return c.setup(c.httpReporter)
}

func (c *Client) setup(zkReporter *zgrp.Reporter) error {
	// 创建一个服务节点, 作为链路上的一个节点
	endpoint, err := zipkin.NewEndpoint(c.serviceName, c.serviceHost+":"+c.servicePort)
	if err != nil {
		return fmt.Errorf("创建节点出错: %v", err)
	}

	nativeTracer, err := zipkin.NewTracer(*zkReporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("创建追踪器错误: %v", err)
	}

	// 把追踪器包装成opentracing
	t := zipkinot.Wrap(nativeTracer)
	c.tracer = &t
	return nil
}

func (c *Client) CloseHttp() {
	if c.httpReporter != nil {
		(*c.httpReporter).Close()
	}
}

func (c *Client) CloseKafka() {
	if c.kafkaReporter != nil {
		(*c.kafkaReporter).Close()
	}
}

func (c *Client) Close() {
	c.CloseHttp()
	c.CloseKafka()
}
