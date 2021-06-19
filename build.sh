#!/bin/bash
echo "start build task"
root_dir=$(pwd)
candidate_dir=$root_dir/server/candidate
pay_dir=$root_dir/server/pay
search_dir=$root_dir/server/search
ticket_dir=$root_dir/server/ticket
ticketPool_dir=$root_dir/server/ticketPool
user_dir=$root_dir/server/user
echo "building candidate server"
cd $candidate_dir
CGO_ENABLED=0 go build -o candidate .
echo "building pay server"
cd $pay_dir
CGO_ENABLED=0 go build -o pay .
echo "building search server"
cd $search_dir
CGO_ENABLED=0 go build -o search .
echo "building ticket server"
cd $ticket_dir
CGO_ENABLED=0 go build -o ticket .
echo "building ticketPool server"
cd $ticketPool_dir
CGO_ENABLED=0 go build -o ticket-pool-k .
echo "building user server"
cd $user_dir
CGO_ENABLED=0 go build -o user .
echo "build finish"