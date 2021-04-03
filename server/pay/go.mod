module pay

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/frankban/quicktest v1.4.1 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/go-redis/redis/v8 v8.6.0
	github.com/google/btree v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.2.0 // indirect
	github.com/pascaldekloe/goe v0.1.0 // indirect
	github.com/pierrec/lz4 v2.2.6+incompatible // indirect
	github.com/segmentio/kafka-go v0.4.10
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb
	google.golang.org/grpc v1.35.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	pay => ../../server/pay/
	rpc => ../../rpc/
)
