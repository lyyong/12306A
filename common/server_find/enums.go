// @Author LiuYong
// @Created at 2021-02-01
// @Modified at 2021-02-01
package server_find

const (
	SCHEME = "consul" //只是一个标识,让grpc.Dial()执行时可以找到该Builder, 然后通过该builder创建resolver
)
