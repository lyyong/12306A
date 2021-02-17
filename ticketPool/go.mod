module ticketPool

go 1.15

require (
	12306A/ticketPool v0.0.0-00010101000000-000000000000
	common v0.0.0-00010101000000-000000000000
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/go-ini/ini v1.62.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/mozillazg/go-pinyin v0.18.0
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	12306A/ticketPool => ../ticketPool
	common => ../common
	rpc => ../rpc
)
