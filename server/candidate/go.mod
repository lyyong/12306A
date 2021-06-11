module candidate

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/go-redis/redis/v8 v8.9.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	pay v0.0.0-00010101000000-000000000000
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	pay => ../pay
	rpc => ../../rpc/
)
