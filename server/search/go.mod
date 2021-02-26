module 12306A-search

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	rpc => ../../rpc
)
