/*
* @Author: 余添能
* @Date:   2021/3/3 8:30 下午
 */
package settings

import (
	"common/tools/logging"
	"testing"
)

func TestSetting(t *testing.T) {
	logging.Setup()
	//*configPath = "../../config/"
	Setup()
	t.Log(RedisDB)
	t.Log(DB)
}
