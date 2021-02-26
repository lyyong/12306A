go env -w CGO_ENABLED=0
go env -w GOOS=linux
go build -a -installsuffix cgo -o ticket .
echo "构建完成"