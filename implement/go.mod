module implement

go 1.15

replace (
	implement => ../implement
	interface => ../interface
	pay => ../server/pay/
	rpc => ../rpc
	common => ../common
)

require (
	interface v0.0.0-00010101000000-000000000000
	rpc v0.0.0-00010101000000-000000000000
)
