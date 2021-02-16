// @Author: KongLingWen
// @Created at 2021/2/16
// @Modified at 2021/2/16

package ticketpool

import (
	"fmt"
	"math/rand"
	pb "rpc/ticketPool/proto/ticketPoolRPC"
	"sync"
	"testing"
	"time"
)

func TestGetTicket_Validity(t *testing.T) {
	initMockData()

	reqCount := 1000
	req := generateGetTicketData(reqCount)

	resp := execBuyTicket(req)
	seatMap := make(map[string]uint64)
	for i := 0; i < len(resp); i++ {
		response := resp[i]
		for _,ticket := range response.Tickets {
			ticketValue := generateRequestValue(int(ticket.StartStationId),int(ticket.DestStationId))
			seat := ticket.CarriageNumber+ticket.SeatNumber
			if ticketValue & seatMap[seat] == 0 {
				seatMap[seat] = seatMap[seat] | ticketValue
			}else {
				t.Fatal("Repeat Ticket")
			}
		}
	}
	fmt.Println("No Repeat Ticket!")
}

func TestGetTicket_Efficient(t *testing.T) {
	//初始化票池
	initMockData()
	// 生成测试数据
	reqCount := 10000
	req := generateGetTicketData(reqCount)

	// 出票并统计耗时
	start := time.Now()
	resp := execBuyTicket(req)
	expend := time.Since(start)

	// 打印出票结果
	printResponse(resp)

	fmt.Println("requestCount:", reqCount, "   time-expend:", expend)
}

func TestSearchTicketCount(t *testing.T) {
	initMockData()
	reqCount := 2000
	getNumberReq := generateGetTicketNumberData(reqCount)
	getTicketReq := generateGetTicketData(reqCount)

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < reqCount; i++ {
		wg.Add(1)
		go func(j int) {
			tp.GetTicket(getTicketReq[j])
			fmt.Println(tp.SearchTicketCount(getNumberReq[j]))	// 打印操作比较耗时，测试效率时删掉输出
			wg.Done()
		}(i)

	}
	wg.Wait()
	expend := time.Since(start)
	fmt.Println("requestCount:", reqCount, "   time-expend:", expend)
}

func TestSearchTicketCount_Validity(t *testing.T) {
	initMockData()
	reqCount := 20
	getTicketReq := generateGetTicketData(reqCount)

	getNumberReq := pb.GetTicketNumberRequest{
		TrainId:        []int32{0},
		StartStationId: 0,
		DestStationId:  5,
		Date:           "2021-02-16",
	}
	for i := 0; i < reqCount; i++ {
		fmt.Println(tp.GetTicket(getTicketReq[i]))
		fmt.Println(tp.SearchTicketCount(getNumberReq))
		fmt.Println()
	}
}

func execBuyTicket(req []pb.GetTicketRequest) []pb.GetTicketResponse {
	resp := make([]pb.GetTicketResponse, len(req))
	var wg sync.WaitGroup
	for i := 0; i < len(req); i++ {
		wg.Add(1)
		go func(j int) {
			// 请求传入票池，出票
			resp[j] = tp.GetTicket(req[j])
			wg.Done()
		}(i)
	}
	wg.Wait()
	return resp
}

func generateGetTicketData(reqCount int) []pb.GetTicketRequest {
	// 请求个数，每个请求包含随机 1~5 张票
	req := make([]pb.GetTicketRequest, reqCount)

	var maxStationNum int32 = 19

	maxPassengerNum := 5
	ticketCount := 0
	for i := 0; i < reqCount; i++ {
		// 随机 1~5 个乘客
		passengerCount := rand.Intn(maxPassengerNum) + 1
		ticketCount += passengerCount
		passengers := make([]*pb.PassengerInfo, passengerCount)
		for j := 0; j < passengerCount; j++ {
			// 随机座位类型，忽略 passengerId 和选座 字段
			passengers[j] = &pb.PassengerInfo{
				PassengerId: 0,
				SeatTypeId:  rand.Int31n(3), // 座位类型 id 随机0-2之间[0,3) [0:商务座，1:一等座，2:二等座]
				ChooseSeat:  "",
			}
		}

		destStation := rand.Int31n(maxStationNum)+1 	// 0 < destStationId <= maxStationNum
		req[i] = pb.GetTicketRequest{
			TrainId:        0,
			StartStationId: rand.Int31n(destStation),	// 0 <= startStationId < destStationId 	----- [0,destStation)
			DestStationId:  destStation,
			Date:           "2021-02-16",
			Passengers:     passengers,
		}
	}
	return req
}

func generateGetTicketNumberData(requestCount int) []pb.GetTicketNumberRequest{
	var maxStationNum int32 = 19

	req := make([]pb.GetTicketNumberRequest, requestCount)

	for i := 0; i < requestCount; i++ {
		destStation := rand.Int31n(maxStationNum)+1
		req[i] = pb.GetTicketNumberRequest{
			TrainId:        []int32{0},
			StartStationId: rand.Int31n(destStation),
			DestStationId:  destStation,
			Date:           "2021-02-16",
		}
	}

	return req
}

func printResponse(resp []pb.GetTicketResponse){

	for i := 0; i < len(resp); i++ {
		response := resp[i]
		fmt.Printf("response %d:\n", i)
		for _,ticket := range response.Tickets {
			fmt.Println(ticket)
		}
	}

}

func TestInit(t *testing.T) {
	initMockData()
	for key, value := range tp.carriageSeatInfoMap {
		fmt.Println("carriageTypeId:[", key, "]; seatInfo:[", value, "]")
	}
	for key, value := range tp.trainMap {
		fmt.Println("trainId:[", key, "]; train:[", value, "]")
		train := value
		for stationId, station := range train.stopStationMap {
			fmt.Println("stationId:[", stationId, "]; station:[", station, "]")
		}
		for date, carriages := range train.carriageMap {
			fmt.Println("date:[", date, "]; carriages:[", carriages, "]")
			for seatTypeId, csi := range carriages.carriageSeatInfo{
				fmt.Println("seatTypeId:[", seatTypeId, "]; csi:[", csi, "]")
				for _,v := range csi.fullTickets {
					fmt.Println("carriage:", v)
				}
			}
		}
	}
}


