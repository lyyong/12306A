// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package database

import "testing"

func TestConnect(t *testing.T) {
	err := Setup("mysql", "root", "12306A.12306A", "localhost:3310", "12306a_test")
	if err != nil {
		t.Error(err)
	}
	Close()
}
