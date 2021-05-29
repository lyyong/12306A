/*
* @Author: 余添能
* @Date:   2021/2/25 2:01 上午
 */
package dao

import (
	"fmt"
	"testing"
)

func TestSelectTrainPoolAll(t *testing.T) {
	trainPools := SelectTrainPoolAll()
	for _, train := range trainPools {
		fmt.Println(train)
	}
}
