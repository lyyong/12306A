购票服务
功能：
    购票
    退票
    改签
    

内部 RPC 接口：
// 添加车票到 db 
rpc AddTickets (Tickets) returns (Empty){}
// 通过订单id 查询车票 
rpc GetTicketByIndentId (GetByIndentRequest) returns (Tickets){}
// 通过乘客id 查询车票
rpc GetTicketByPassengerId (GetByPassengerRequest) returns (Tickets){}
// 更新车票状态 （退票/改签...）
rpc UpdateState (UpdateStateRequest) returns (Empty){}

外部接口：
router:"/buyTicket" 
