module implement

go 1.15

replace (
	common => ../common
	implement => ../implement
	interface => ../interface
	pay => ../server/pay/
	rpc => ../rpc
)

require (
	common v0.0.0-00010101000000-000000000000
	interface v0.0.0-00010101000000-000000000000
	rpc v0.0.0-00010101000000-000000000000
)
