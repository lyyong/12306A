module rpc

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3
	github.com/segmentio/kafka-go v0.4.10
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
)

replace (
	common => ../common
	rpc => ../rpc
)
