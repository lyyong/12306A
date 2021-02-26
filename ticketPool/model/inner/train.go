/*
* @Author: 余添能
* @Date:   2021/2/25 11:43 下午
 */
package inner

type Train struct {
	ID int

	CreatedAt string
	UpdatedAt string
	DeletedAt string
	CreatedBy string
	UpdatedBy string
	DeletedBy string

	Number string
	StartStation string
	EndStation string
	TrainType int
	State int
}
