module user

go 1.15

require (
	common v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/spf13/cobra v1.1.1
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.5.1
	gopkg.in/ini.v1 v1.51.0
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.20.12
	rpc v0.0.0-00010101000000-000000000000
)

replace (
	common => ../../common
	rpc => ../../rpc
)
