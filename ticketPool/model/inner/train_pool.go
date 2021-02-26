/*
* @Author: 余添能
* @Date:   2021/2/25 1:57 上午
 */
package inner

//通过坐车时间和城市去找合适的车次
type TrainPool struct {
	ID int
	CreatedAt string
	UpdatedAt string
	DeletedAt string
	CreatedBy string
	UpdatedBy string
	DeletedBy string

	TrainNo   string
	StartCity string
	//StartStation string
	StartTime string
	EndCity   string
	//EndStation   string
	EndTime string
}
