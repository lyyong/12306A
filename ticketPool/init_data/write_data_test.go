/*
* @Author: 余添能
* @Date:   2021/1/23 1:34 上午
 */
package init_data

import (
	"fmt"
	"testing"
)

func TestWriteStationProvinceCity(t *testing.T) {
	WriteStationProvinceCity()
}

func TestReadCity(t *testing.T) {
	ReadCity()
}

func TestReadTrainNo(t *testing.T) {
	trains := ReadTrainNo()
	fmt.Println(len(trains))
}

func TestWriteTotalTrainNo(t *testing.T) {
	WriteTotalTrainNo()
}
func TestReadStationCity(t *testing.T) {
	stationCity := ReadStationCity()
	fmt.Println(stationCity["明光"])
}
