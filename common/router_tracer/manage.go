// @Author LiuYong
// @Created at 2021-02-01
// @Modified at 2021-02-01
package router_tracer

import (
	"common/tools/logging"
	"errors"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var zipkinClient *Client

// 懒汉式双重检查单例
// var lock sync.Once

// SetupByHttp 开启链路追踪功能, 通过http向zipkin服务器进行汇报
// serviceName - 当前服务的名称
// serviceHost - 当前服务的host 不包括端口
// servicePort - 当前服务的端口
// zkHttp - zipkin服务器汇报的http api,例如"http://localhost:9411/api/v2/spans"
func SetupByHttp(serviceName, serviceHost, servicePort, zkHttp string) error {
	// 创建一个报告器, 用来想zipkin汇报信息
	if zipkinClient != nil {
		return errors.New("重复创建zipkin-http")
	}
	logging.Info("开启链路追踪功能, 当前服务的名称: ", serviceName, " 当前服务Host: ", serviceHost, " 当前服务的端口号: ", servicePort, " zipkin的http地址: "+zkHttp)
	zipkinClient = newClient(serviceName, serviceHost, servicePort)
	rp := zipkinhttp.NewReporter(zkHttp)
	zipkinClient.httpReporter = &rp

	if err := zipkinClient.setup(zipkinClient.httpReporter); err != nil {
		zipkinClient = nil
		logging.Error("链路追踪功能开启失败")
		return err
	}
	logging.Info("链路追踪功能开启成功")
	return nil
}

// GetClient 获取 router_tracer.Client 单例
func GetClient() (*Client, error) {
	if zipkinClient == nil {
		return nil, errors.New("未使用zipkin.SetupByHttp()创建zipkin.Client实例")
	}
	return zipkinClient, nil
}

// IsTracing 是否有单例存在
func IsTracing() bool {
	return zipkinClient != nil
}

// Close 关闭 router_tracer.Client
func Close() {
	zipkinClient.closeHttp()
	zipkinClient.closeKafka()
	zipkinClient = nil
}
