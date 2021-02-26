/*
* @Author: 余添能
* @Date:   2021/1/31 6:37 下午
 */
package init_redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"ticketPool/dao"
	"ticketPool/rdb"
	"time"
)


func InitDataRedis()  {

	//初始化trainpool
	fmt.Println("初始化redis：trains")
	WriteTrainPoolToRedis()
	//车站
	//fmt.Println("初始化redis：station")
	//WriteStationToRedis()
	//车站:城市
	fmt.Println("初始化redis：station-city映射")
	WriteStationAndCityToRedis()
	//列车信息
	fmt.Println("开始初始化redis：列车基本信息")
	WriteTrainInfoToRedis()
	//票池
	//fmt.Println("开始初始化redis：票池")
	//WriteTicketPoolToRedis()
	//城市之间的车次


}


//用出发城市：目的城市作为key，使用zset类型存储两个城市之间的车次
func WriteTrainPoolToRedis() {
	trainPools := dao.SelectTrainPoolAll()
	for _, trainPool := range trainPools {
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", trainPool.StartTime, time.Local)
		//endTime,_:=time.ParseInLocation("2006-01-02 15:04:05", trainPool.EndTime, time.Local)
		//以上车时间作为score,小时+分钟，不用日期，因为车次每天都会有
		hour, minute, _ := startTime.Clock()
		start := hour*60 + minute
		//fmt.Println(start)
		key := trainPool.StartCity + "-" + trainPool.EndCity

		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(start), Member: trainPool.TrainNo})
		//fmt.Println(cmd.Result())

	}
	//rdb.RedisDB.ZRangeByScoreWithScores("景德镇-九江", redis.ZRangeBy{Min: "0", Max: "50000"})
	//fmt.Println(zslice.Val())
}

//假设所有车次天天会有，所以查车次不用日期
//但可能出现停运情况,特殊考虑


//使用hash存储列车基本信息,key: trainNo 如, G104
//列车信息：列车基本信息、车站信息
//列车基本信息用一个hash保存：key=车次, 元素：trainNo,stationNum, 1--stationNum -> cityName
//每个车站都单独用一个hash保存,key=车次+站序，元素：stationNo,stationName,cityName,arriveTime,departTime,duration,price,mileage
func WriteTrainInfoToRedis() {
	stopInfos:=dao.SelectStopInfoAll()
	for _,stopInfo:=range stopInfos{
		key:=stopInfo.TrainNumber
		rdb.RedisDB.HSet(key,"stationNum",0)
	}

	for _,stopInfo:=range stopInfos{
		key:=stopInfo.TrainNumber
		//车次元数据
		//站名-站序
		rdb.RedisDB.HSet(key,stopInfo.StationName,stopInfo.StopSeq)
		//城市-站名
		rdb.RedisDB.HSet(key,stopInfo.City,stopInfo.StationName)
		//记录站数
		rdb.RedisDB.HIncrBy(key,"stationNum",1)
		//rdb.RedisDB.HSet(key,stopInfo.City,stopInfo.StationName)

		//保存站点信息:G21-1
		stationKey := stopInfo.TrainNumber + "-" + strconv.Itoa(stopInfo.StopSeq)
		rdb.RedisDB.HSet(stationKey, "stopSeq", stopInfo.StopSeq)
		rdb.RedisDB.HSet(stationKey, "stationName", stopInfo.StationName)
		rdb.RedisDB.HSet(stationKey, "cityName", stopInfo.City)
		//到达时间
		//fmt.Println(stopInfo.ArrivedTime,stopInfo.LeaveTime)
		rdb.RedisDB.HSet(stationKey, "arriveTime", stopInfo.ArrivedTime)
		//出发时间
		rdb.RedisDB.HSet(stationKey, "leaveTime", stopInfo.LeaveTime)
		//保存持续时间
		//rdb.RedisDB.HSet(stationKey, "duration", stopInfo.StayDuration)
		//rdb.RedisDB.HSet(stationKey, "mileage", stopInfo.Mileage)
	}
}

//用来一个hash保存 站点-城市 的映射，
func WriteStationAndCityToRedis() {
	stations:=dao.SelectStationAll()
	if stations==nil{
		return
	}
	key := "stationCity"
	for _,s:=range stations {
		rdb.RedisDB.HSet(key, s.Name, s.City)
	}
}

//上海虹桥

//用一个List保存所有站点
func WriteStationToRedis()  {
	stations:=dao.SelectStationAll()

	key := "stations"
	for _,v:=range stations {
		spell:=v.Spell[0]
		//fmt.Println(v,spell)
		rdb.RedisDB.ZAdd(key,redis.Z{Score: float64(spell),Member: v.Name})
		//rdb.RedisDB.LPush(key,v.Name)
	}
}




//初始化train_pools，用于查找两个城市之间的车次
//只写入北京-上海的车次数据
//func WriteTrainPool() {
//
//	res,err:=rdb.RedisDB.ZRangeByScore("北京-上海",redis.ZRangeBy{Min: "0",Max: "10000"}).Result()
//	if err!=nil{
//		fmt.Println(err)
//		return
//	}
//	var resMap map[string]string
//	resMap=make(map[string]string)
//	for _,v:=range res{
//		resMap[v]="111"
//		//fmt.Println(v)
//	}
//	fmt.Println(res)
//	fmt.Println("开始初始化train_pool表")
//	trains := init_data.ReadTotalTrainNo()
//
//	sqlStr := "insert into train_pools(created_at,created_by,initial_time,terminal_time,train_number,start_city,start_time,end_city,end_time) " +
//		"values(?,?,?,?,?,?,?,?,?);"
//	st, err := init_data.Db.Prepare(sqlStr)
//	defer st.Close()
//	if err != nil {
//		fmt.Println("prepare table train_pools failed, err:", err)
//		return
//	}
//
//	for _, train := range trains {
//
//		if resMap[train.TrainNo]==""{
//			continue
//		}
//		n := len(train.Stations)
//		stations := train.Stations
//		//记录上车城市是否写入
//		startCityWrited := make(map[string]string)
//		for i := 0; i < n; i++ {
//			//上车城市
//			startCity := stations[i].CityName
//			if startCityWrited[startCity] != "" {
//				continue
//			}
//			startCityWrited[startCity] = startCity
//			//记录下车城市是否写入
//			endCityWrited := make(map[string]string)
//			for j := i + 1; j < n; j++ {
//				//下车城市
//				endCity := stations[j].CityName
//				//去重
//				//一趟车可能会经过同一个城市的两个车站，比如下属县级市，不重复输入
//				if endCityWrited[endCity] == "" {
//					//车次的：起始时间，车次终止时间，车次，上车城市，上车时间，下车城市，下车时间
//					st.Exec(time.Now(),"系统",train.InitialTime, train.TerminalTime, train.TrainNo, startCity, stations[i].ArriveTime, endCity, stations[j].ArriveTime)
//					endCityWrited[endCity] = endCity
//				}
//			}
//		}
//	}
//}

