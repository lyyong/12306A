/*
* @Author: 余添能
* @Date:   2021/2/21 11:33 下午
 */
package init_redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"ticketPool/dao"
	"ticketPool/rdb"
)

func InitTicketPool()  {

}

//北京市-上海市 0 1000
//将mysql中carriage_type数据读入，5种车厢，依次循环

//初始化每趟车次的票池,全票
//second,first,business
//zset保存票
//key=日期:车次::
//date:="2021-02-25"
func WriteTicketPoolToRedis()  {
	////上海-北京的车次
	trainNos,err:=rdb.RedisDB.ZRangeByScore("北京市-上海市",redis.ZRangeBy{Min: "0",Max: "1000"}).Result()
	if err!=nil{
		fmt.Println("ZRangeByScore failed, err:",err)
		return
	}
	//trainNos:=dao.SelectTrainAll()
	fmt.Println(len(trainNos))
	carriageTypes:=dao.QueryCarriageTypesAll()
	carriageTypeNum:=len(carriageTypes)
	carriageNum:=30
	for i:=0;i<len(trainNos);i++ {
		trainNo := trainNos[i]
		//fmt.Println(trainNo)
		trainMap,_:=rdb.RedisDB.HGetAll(trainNo).Result()
		stationNum,_:=strconv.ParseInt(trainMap["stationNum"],10,64)
		//循环
		k := 0
		for j := 0; j < carriageNum; j++ {
			carriageType := carriageTypes[k]
			k++
			if k==carriageTypeNum{
				k=0
			}
			//stationKey:=trainNo+"-"+"1"
			//stationMap,_:=rdb.RedisDB.HGetAll(stationKey).Result()
			//depart,_:=time.Parse("2006-01-02 15:04:05",stationMap["departTime"])
			//departTime:=depart.Add(time.Now().Sub(depart)).Format("2006-01-02")
			//fmt.Println(depart,carriageType,stationNum)
			date:="2021-02-25"
			if carriageType.BusinessSeatNumber != 0 {
				seats := strings.Split(carriageType.BusinessSeat, ",")
				for _, seat := range seats {

					key := date + ":" + trainNo + ":" + "1:" + "businessSeat"
					//member=车厢号:座位号
					rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(stationNum), Member: strconv.Itoa(j+1) + ":" + seat})
				}
			}
			if carriageType.FirstSeatNumber != 0 {
				seats := strings.Split(carriageType.FirstSeat, ",")
				for _, seat := range seats {
					key := date + ":" + trainNo + ":" + "1:" + "firstSeat"
					//member=车厢号:座位号
					rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(stationNum), Member: strconv.Itoa(j+1) + ":" + seat})
				}
			}
			if carriageType.SecondSeatNumber != 0 {
				seats := strings.Split(carriageType.SecondSeat, ",")
				for _, seat := range seats {
					key := date + ":" + trainNo + ":" + "1:" + "secondSeat"
					//member=车厢号:座位号
					rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(stationNum), Member: strconv.Itoa(j+1) + ":" + seat})
				}
			}
		}
	}
}
//将票分散到各个站点，均分

func SplitTicket()  {
	////上海-北京的车次
	trainNos,err:=rdb.RedisDB.ZRangeByScore("北京市-上海市",redis.ZRangeBy{Min: "0",Max: "3000"}).Result()
	if err!=nil{
		fmt.Println("ZRangeByScore failed, err:",err)
		return
	}
	//trainNos:=dao.SelectTrainAll()
	fmt.Println(len(trainNos))
	for i:=0;i<len(trainNos);i++ {
		trainNo := trainNos[i]
		//fmt.Println(trainNo)
		trainMap, _ := rdb.RedisDB.HGetAll(trainNo).Result()
		stationNum, _ := strconv.ParseInt(trainMap["stationNum"], 10, 64)
		date := "2021-02-25"
		//一等座
		firstSeatKey := date + ":" + trainNo + ":" + "1" +":"+ "firstSeat"
		firstSeats, _ := rdb.RedisDB.ZRangeByScoreWithScores(firstSeatKey, redis.ZRangeBy{Min: "0", Max: "1000"}).Result()
		firstSeatMean := len(firstSeats) / int(stationNum)
		//fmt.Println(trainNo,firstSeatMean,len(firstSeats),stationNum)
		for j:=2;j<int(stationNum);j++{
			for k:=0;k<firstSeatMean;k++{
				tickets,_:=rdb.RedisDB.ZRangeByScoreWithScores(firstSeatKey,redis.ZRangeBy{Min: strconv.Itoa(j+1),Max: "10000"}).Result()
				if tickets==nil{
					break
				}
				//fmt.Println(len(tickets))
				//删掉全票
				rdb.RedisDB.ZRem(firstSeatKey,tickets[0].Member).Result()
				//fmt.Println(t,err)
				//分票并写回
				//fmt.Println(trainNo,j,tickets[0])
				rdb.RedisDB.ZAdd(firstSeatKey,redis.Z{Score: float64(j),Member: tickets[0].Member})
				remainderKey:=date + ":" + trainNo + ":" + strconv.Itoa(j) +":"+ "firstSeat"
				rdb.RedisDB.ZAdd(remainderKey,redis.Z{Score: tickets[0].Score,Member: tickets[0].Member})
			}
		}

		//二等座
		businessSeatKey := date + ":" + trainNo + ":" + "1" +":"+ "businessSeat"
		businessSeats, _ := rdb.RedisDB.ZRangeByScoreWithScores(businessSeatKey, redis.ZRangeBy{Min: "0", Max: "1000"}).Result()
		businessSeatMean := len(businessSeats) / int(stationNum)
		//fmt.Println(trainNo,businessSeatMean,len(businessSeats),stationNum)
		for j:=2;j<int(stationNum);j++{
			for k:=0;k<businessSeatMean;k++{
				tickets,_:=rdb.RedisDB.ZRangeByScoreWithScores(businessSeatKey,redis.ZRangeBy{Min: strconv.Itoa(j+1),Max: "10000"}).Result()
				if tickets==nil{
					break
				}
				//fmt.Println(len(tickets))
				//删掉全票
				rdb.RedisDB.ZRem(businessSeatKey,tickets[0].Member).Result()
				//fmt.Println(t,err)
				//分票并写回
				//fmt.Println(trainNo,j,tickets[0])
				rdb.RedisDB.ZAdd(businessSeatKey,redis.Z{Score: float64(j),Member: tickets[0].Member})
				remainderKey:=date + ":" + trainNo + ":" + strconv.Itoa(j) +":"+ "businessSeat"
				rdb.RedisDB.ZAdd(remainderKey,redis.Z{Score: tickets[0].Score,Member: tickets[0].Member})
			}
		}

		//商务座
		secondSeatKey := date + ":" + trainNo + ":" + "1" +":"+ "secondSeat"
		secondSeats, _ := rdb.RedisDB.ZRangeByScoreWithScores(secondSeatKey, redis.ZRangeBy{Min: "0", Max: "1000"}).Result()
		secondSeatMean := len(secondSeats) / int(stationNum)
		fmt.Println(trainNo,secondSeatMean,len(secondSeats),stationNum)
		for j:=2;j<int(stationNum);j++{
			for k:=0;k<secondSeatMean;k++{
				tickets,_:=rdb.RedisDB.ZRangeByScoreWithScores(secondSeatKey,redis.ZRangeBy{Min: strconv.Itoa(j+1),Max: "10000"}).Result()
				if tickets==nil{
					break
				}
				//fmt.Println(len(tickets))
				//删掉全票
				rdb.RedisDB.ZRem(secondSeatKey,tickets[0].Member).Result()
				//fmt.Println(t,err)
				//分票并写回
				//fmt.Println(trainNo,j,tickets[0])
				rdb.RedisDB.ZAdd(secondSeatKey,redis.Z{Score: float64(j),Member: tickets[0].Member})
				remainderKey:=date + ":" + trainNo + ":" + strconv.Itoa(j) +":"+ "secondSeat"
				rdb.RedisDB.ZAdd(remainderKey,redis.Z{Score: tickets[0].Score,Member: tickets[0].Member})
			}
		}
	}
}

//if carriageType.SeniorSoftBenthNumber != 0 {
//	seats := strings.Split(carriageType.SeniorSoftBenth, ",")
//	for _, seat := range seats {
//		departTime := train.Stations[0].DepartTime.Format("2006-01-02")
//		key := departTime + ":" + train.TrainNo + ":" + "1:" + "seniorSoftBench"
//		//member=车厢号:座位号
//		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(train.StationNum), Member: strconv.Itoa(j+1) + ":" + seat})
//	}
//}
//if carriageType.SoftBenthNumber != 0 {
//	seats := strings.Split(carriageType.SoftBench, ",")
//	for _, seat := range seats {
//		departTime := train.Stations[0].DepartTime.Format("2006-01-02")
//		key := departTime + ":" + train.TrainNo + ":" + "1:" + "softBench"
//		//member=车厢号:座位号
//		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(train.StationNum), Member: strconv.Itoa(j+1) + ":" + seat})
//	}
//}
//if carriageType.HardBenthNumber != 0 {
//	seats := strings.Split(carriageType.HardBenth, ",")
//	for _, seat := range seats {
//		departTime := train.Stations[0].DepartTime.Format("2006-01-02")
//		key := departTime + ":" + train.TrainNo + ":" + "1:" + "hardBench"
//		//member=车厢号:座位号
//		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(train.StationNum), Member: strconv.Itoa(j+1) + ":" + seat})
//	}
//}
//if carriageType.HardSeatNumber != 0 {
//	seats := strings.Split(carriageType.HardSeat, ",")
//	for _, seat := range seats {
//		departTime := train.Stations[0].DepartTime.Format("2006-01-02")
//		key := departTime + ":" + train.TrainNo + ":" + "1:" + "hardSeat"
//		//member=车厢号:座位号
//		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(train.StationNum), Member: strconv.Itoa(j+1) + ":" + seat})
//	}
//}


//北京-上海 0 1000
//
////将票分散到各个站
//func InitTicketPool()  {
//	trainNos:= ReadTrainNoFromSToB()
//
//	for _,train:=range trainNos{
//
//		//trainNo:=train.TrainNo
//		//获取车次元数据
//		//resMap,_:= rdb.RedisDB.HGetAll(trainNo).Result()
//		stationNum:=train.StationNum
//		departTime:=train.Stations[0].DepartTime.Format("2006-01-02")
//
//		businessSeatKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "businessSeat"
//		businessSeats,_:= rdb.RedisDB.ZRangeByScoreWithScores(businessSeatKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		businessSeatNum:=len(businessSeats)
//
//
//		businessSeatMean:=businessSeatNum/int(stationNum)
//		for i:=2;i<int(stationNum);i++{
//			businessSeats,_= rdb.RedisDB.ZRangeByScoreWithScores(businessSeatKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//			for j:=0;j<businessSeatMean;j++{
//				rdb.RedisDB.ZRem(businessSeatKey,businessSeats[j].Member)
//
//			}
//		}
//		//商务座
//		for index:=0;index<businessSeatNum;index++{
//			startStationNo:=index/int(stationNum)+1
//			//拆票
//			ticket:=&outer.Ticket{}
//			ticket.Date=departTime
//			ticket.TrainNo=train.TrainNo
//			ticket.StartTime=train.Stations[startStationNo-1].DepartTime.Format("2006-01-02 15:04:05")
//			ticket.StartStation=train.Stations[startStationNo-1].StationName
//			ticket.StartStationNo=strconv.Itoa(startStationNo)
//			ticket.EndTime=train.Stations[stationNum-1].ArriveTime.Format("2006-01-02 15:04:05")
//			ticket.EndStation=train.Stations[stationNum-1].StationName
//			ticket.EndStationNo=strconv.Itoa(int(train.Stations[stationNum-1].StationNo))
//
//			carriageAndSeatNo := strings.Split(businessSeats[index].Member.(string), ":")
//			carriageNo := carriageAndSeatNo[0]
//			seatNo := carriageAndSeatNo[1]
//			ticket.CarriageNo = carriageNo
//			ticket.SeatNo = seatNo
//			dao.InsertTicket(ticket)
//			if startStationNo>1{
//				//余票
//				remainder:=&outer.Ticket{}
//				remainder.Date=departTime
//				remainder.TrainNo=train.TrainNo
//				remainder.StartTime=train.Stations[startStationNo-1].DepartTime.Format("2006-01-02 15:04:05")
//				remainder.StartStation=train.Stations[startStationNo-1].StationName
//				remainder.StartStationNo=strconv.Itoa(1)
//				remainder.EndTime=train.Stations[stationNum-1].ArriveTime.Format("2006-01-02 15:04:05")
//				remainder.EndStation=train.Stations[stationNum-1].StationName
//				remainder.EndStationNo=strconv.Itoa(startStationNo)
//				//carriageAndSeatNo := strings.Split(businessSeats[index].Member.(string), ":")
//				//carriageNo := carriageAndSeatNo[0]
//				//seatNo := carriageAndSeatNo[1]
//				remainder.CarriageNo = carriageNo
//				remainder.SeatNo = seatNo
//			}
//		}
//		//for i:=2;i<=int(stationNum);i++{
//		//	for j:=0;j<businessSeatMean;j++{
//		//		//拆票
//		//		ticket:=&outer.Ticket{}
//		//		ticket.Date=departTime
//		//		ticket.TrainNo=train.TrainNo
//		//		ticket.StartTime=train.Stations[i].DepartTime.Format("2006-01-02 15:04:05")
//		//		ticket.StartStation=train.Stations[i].StationName
//		//		ticket.StartStationNo=i
//		//		ticket.EndTime=train.Stations[stationNum-1].ArriveTime.Format("2006-01-02 15:04:05")
//		//		ticket.EndStation=train.Stations[stationNum-1].StationName
//		//		ticket.EndStationNo=int(train.Stations[stationNum-1].StationNo)
//		//
//		//		carriageAndSeatNo := strings.Split(businessSeats[index].Member.(string), ":")
//		//		carriageNo := carriageAndSeatNo[0]
//		//		seatNo := carriageAndSeatNo[1]
//		//		ticket.CarriageNo = carriageNo
//		//		ticket.SeatNo = seatNo
//		//		dao.InsertTicket(ticket)
//		//		index++
//		//	}
//		//}
//		////第一站的票
//		//for ;index<businessSeatNum;index++{
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: float64(1),Member: businessSeats[index].Member})
//		//}
//		//
//		//firstSeatKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "firstSeat"
//		//firstSeats,_:=RedisDB.ZRangeByScoreWithScores(firstSeatKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		//firstSeatNum:=len(firstSeats)
//		//for i:=2;i<=int(stationNum);i++{
//		//	businessSeatMean:=businessSeatNum/int(stationNum)
//		//	firstSeatMean:=firstSeatNum/int(stationNum)
//		//	secondSeatMean:=secondSeatNum/int(stationNum)
//		//	seniorSoftBenchMean:=seniorSoftBenchNum/int(stationNum)
//		//	softBenchMean:=softBenchNum/int(stationNum)
//		//	hardBenchMean:=hardBenchNum/int(stationNum)
//		//	hardSeatMean:=hardSeatNum/int(stationNum)
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: i})
//		//}
//		//
//		//secondSeatKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "firstSeat"
//		//secondSeats,_:=RedisDB.ZRangeByScoreWithScores(secondSeatKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		//secondSeatNum:=len(secondSeats)
//		//for i:=2;i<=int(stationNum);i++{
//		//	businessSeatMean:=businessSeatNum/int(stationNum)
//		//	firstSeatMean:=firstSeatNum/int(stationNum)
//		//	secondSeatMean:=secondSeatNum/int(stationNum)
//		//	seniorSoftBenchMean:=seniorSoftBenchNum/int(stationNum)
//		//	softBenchMean:=softBenchNum/int(stationNum)
//		//	hardBenchMean:=hardBenchNum/int(stationNum)
//		//	hardSeatMean:=hardSeatNum/int(stationNum)
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: i})
//		//}
//		//
//		//seniorSoftBenchKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "firstSeat"
//		//seniorSoftBenchs,_:=RedisDB.ZRangeByScoreWithScores(seniorSoftBenchKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		//seniorSoftBenchNum:=len(seniorSoftBenchs)
//		//for i:=2;i<=int(stationNum);i++{
//		//	businessSeatMean:=businessSeatNum/int(stationNum)
//		//	firstSeatMean:=firstSeatNum/int(stationNum)
//		//	secondSeatMean:=secondSeatNum/int(stationNum)
//		//	seniorSoftBenchMean:=seniorSoftBenchNum/int(stationNum)
//		//	softBenchMean:=softBenchNum/int(stationNum)
//		//	hardBenchMean:=hardBenchNum/int(stationNum)
//		//	hardSeatMean:=hardSeatNum/int(stationNum)
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: i})
//		//}
//		//
//		//softBenchKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "firstSeat"
//		//softBenchs,_:=RedisDB.ZRangeByScoreWithScores(softBenchKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		//softBenchNum:=len(softBenchs)
//		//for i:=2;i<=int(stationNum);i++{
//		//	businessSeatMean:=businessSeatNum/int(stationNum)
//		//	firstSeatMean:=firstSeatNum/int(stationNum)
//		//	secondSeatMean:=secondSeatNum/int(stationNum)
//		//	seniorSoftBenchMean:=seniorSoftBenchNum/int(stationNum)
//		//	softBenchMean:=softBenchNum/int(stationNum)
//		//	hardBenchMean:=hardBenchNum/int(stationNum)
//		//	hardSeatMean:=hardSeatNum/int(stationNum)
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: i})
//		//}
//		//
//		//hardBenchKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "firstSeat"
//		//hardBenchs,_:=RedisDB.ZRangeByScoreWithScores(hardBenchKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		//hardBenchNum:=len(hardBenchs)
//		//for i:=2;i<=int(stationNum);i++{
//		//	businessSeatMean:=businessSeatNum/int(stationNum)
//		//	firstSeatMean:=firstSeatNum/int(stationNum)
//		//	secondSeatMean:=secondSeatNum/int(stationNum)
//		//	seniorSoftBenchMean:=seniorSoftBenchNum/int(stationNum)
//		//	softBenchMean:=softBenchNum/int(stationNum)
//		//	hardBenchMean:=hardBenchNum/int(stationNum)
//		//	hardSeatMean:=hardSeatNum/int(stationNum)
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: i})
//		//}
//		//
//		//hardSeatKey:=departTime + ":" + train.TrainNo + ":" + "1:" + "firstSeat"
//		//hardSeats,_:=RedisDB.ZRangeByScoreWithScores(hardSeatKey,redis.ZRangeBy{Max: "10000",Min: "0"}).Result()
//		//hardSeatNum:=len(hardSeats)
//		//
//		//for i:=2;i<=int(stationNum);i++{
//		//	businessSeatMean:=businessSeatNum/int(stationNum)
//		//	firstSeatMean:=firstSeatNum/int(stationNum)
//		//	secondSeatMean:=secondSeatNum/int(stationNum)
//		//	seniorSoftBenchMean:=seniorSoftBenchNum/int(stationNum)
//		//	softBenchMean:=softBenchNum/int(stationNum)
//		//	hardBenchMean:=hardBenchNum/int(stationNum)
//		//	hardSeatMean:=hardSeatNum/int(stationNum)
//		//	RedisDB.ZAdd(businessSeatKey,redis.Z{Score: i})
//		//}
//
//	}
//}
