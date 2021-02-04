/*
* @Author: 余添能
* @Date:   2021/2/4 11:42 下午
 */
package init_data

func InitDataMysql()  {
	//初始化station_province_city
	WriteStationProvinceCity()

	//初始化total_train_no表
	WriteTotalTrainNo()
	ReviseTotalTrainNo()

	//初始化train_pool表
	WriteTrainPool()


}
