/*
* @Author: 余添能
* @Date:   2021/2/21 11:47 下午
 */
package init_redis

import (
	"testing"
)



func TestWriteTicketPoolToRedis(t *testing.T) {
	WriteTicketPoolToRedis()
}

func TestSplitTicket(t *testing.T) {
	SplitTicket()
}