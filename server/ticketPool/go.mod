module ticketPool

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/go-ini/ini v1.62.0 // indirect
	github.com/gomodule/redigo v1.8.3 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506 // indirect
	google.golang.org/grpc v1.35.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	gorm.io/driver/mysql v1.0.4 // indirect
	gorm.io/gorm v1.20.12 // indirect
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	rpc => ../../rpc/
)
