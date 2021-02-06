module ticket

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gomodule/redigo v1.8.3
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	google.golang.org/grpc v1.35.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.20.12
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	rpc => ../../rpc/
)
