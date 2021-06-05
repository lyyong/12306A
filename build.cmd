echo start build task
go env -w GOOS=linux
set root_dir=%cd%
set candidate_dir=%root_dir%\server\candidate
set pay_dir=%root_dir%\server\pay
set search_dir=%root_dir%\server\search
set ticket_dir=%root_dir%\server\ticket
set ticketPool_dir=%root_dir%\server\ticketPool
set user_dir=%root_dir%\server\user
echo building candidate server
cd %candidate_dir%
go build -o candidate .
echo building pay server
cd %pay_dir%
go build -o pay .
echo building search server
cd %search_dir%
go build -o search .
echo building ticket server
cd %ticket_dir%
go build -o ticket .
echo building ticketPool server
cd %ticketPool_dir%
go build -o ticket-pool-k .
echo building user server
cd %user_dir%
go build -o user .
echo build finish
go env -w GOOS=windows
pause