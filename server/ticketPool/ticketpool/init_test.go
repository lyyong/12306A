// @Author LiuYong
// @Created at 2021-02-21
// @Modified at 2021-02-21
package ticketpool

import (
	"testing"
	"ticketPool/utils/database"
	"ticketPool/utils/setting"
)

func TestInitTicketPool(t *testing.T) {
	setting.InitSetting()
	database.Setup()
	InitTicketPool()
	database.Close()
}

func BenchmarkName(b *testing.B) {
	setting.InitSetting()
	database.Setup()
	for i := 0; i < 1; i++ {
		InitTicketPool()
	}
	database.Close()
}
