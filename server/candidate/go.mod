module candidate

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/go-redis/redis/v8 v8.9.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	gorm.io/gorm v1.21.11
	pay v0.0.0-00010101000000-000000000000
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	pay => ../pay
	rpc => ../../rpc/
)
