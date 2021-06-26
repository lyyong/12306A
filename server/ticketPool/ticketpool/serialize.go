// @Author: KongLingWen
// @Created at 2021/6/19
// @Modified at 2021/6/19

package ticketpool

import (
	"common/tools/logging"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

const (
	JSONFile = "TicketPoolData.json"
	GOBFile  = "TicketPoolData.gob"
)

func Serialize() {
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			tp := Tp
			TpLock.Lock()
			logging.Info("开始序列化")
			st := time.Now()
			// serializeByJSON(&tp)
			serializeByGOB(&tp)
			logging.Info("序列化结束, 耗时: ", time.Now().Sub(st).Milliseconds())
			TpLock.Unlock()
		}
	}()
}

func UnSerialize(dest interface{}) error {
	// return unSerializeByJSON(dest)
	return unSerializeByGOB(dest)
}

// serializeByJSON 通过JSON进行序列化
func serializeByJSON(value interface{}) {
	res, err := json.Marshal(value)
	if err != nil {
		logging.Error(err)
	}
	err = ioutil.WriteFile(JSONFile, res, 0644)
	if err != nil {
		logging.Error(err)
	}
}

// unSerializeByJSON 从JSON文件恢复, dest必须为指针类型
func unSerializeByJSON(dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		// 不为指针类型
		return errors.New("dest参数应该为指针类型")
	}
	file, err := ioutil.ReadFile(JSONFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, dest); err != nil {
		return err
	}
	return nil
}

// serializeByGOB 通过gob进行序列化
func serializeByGOB(value interface{}) {
	file, err := os.OpenFile(GOBFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.Error(err)
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(value)
	if err != nil {
		logging.Error(err)
		return
	}
}

// unSerializeByGOB 通过gob反序列化
func unSerializeByGOB(dest interface{}) error {
	file, err := os.OpenFile(GOBFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err = dec.Decode(dest)
	if err != nil {
		return err
	}
	return nil
}
