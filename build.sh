#!/bin/bash
echo "start build task"
go env -w GOOS=linux
echo "building candidate server"
go build -o ./server/candidate/candidate ./server/candidate/
echo "building pay server"
go build -o ./server/pay/pay ./server/pay/
echo "building search server"
go build -o ./server/search/search ./server/search/
echo "building ticket server"
go build -o ./server/ticket/ticket ./server/ticket/
echo "building ticketPool server"
go build -o ./server/ticketPool/ticket-pool-k ./server/ticketPool/
echo "building user server"
go build -o ./server/user/user ./server/user/
echo "build finish"
echo "enter any button to close window"
pause