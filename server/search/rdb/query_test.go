/*
* @Author: 余添能
* @Date:   2021/2/3 7:30 下午
 */
package rdb

import (
	"12306A/server/search/model/outer"
	"testing"
)

func TestQuery(t *testing.T) {
	search:=&outer.Search{}
	search.StartCity="北京"
	search.EndCity="上海"
	search.Date="2021-1-23 01:02:03"
	//search.SeatClass=""
	Query(search)
}
