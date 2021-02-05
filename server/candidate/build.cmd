swag init -g router/router.go
go env -w CGO_ENABLED=0
go env -w GOOS=linux
go build -a -installsuffix cgo -o candidate .
echo "构建完成"