FROM golang:latest

WORKDIR /home
COPY . .

RUN go env -w GO111MODULE=auto
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /home/server/user
RUN go build -o user main.go

EXPOSE 8222
ENTRYPOINT ["./user", "server"]