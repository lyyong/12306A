module rpc

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	pay v0.0.0-00010101000000-000000000000
)

replace (
	common => ../common
	pay => ../server/pay/
	rpc => ../rpc
)
