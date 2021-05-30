/*
* @Author: 余添能
* @Date:   2021/2/25 1:44 上午
 */
package inner

type StopInfo struct {
	ID int64
	//created_at,modified_at,deleted_at,created_by,modified_by,deleted_by,train_id,station_id,
	CreatedAt string
	UpdatedAt string
	DeletedAt string
	CreatedBy string
	UpdatedBy string
	DeletedBy string

	TrainId     int
	StationId   int
	TrainNumber string
	StationName string
	City        string
	ArrivedTime string
	LeaveTime   string
	//StayDuration int
	StopSeq int
	//Mileage int
}
