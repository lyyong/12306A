// @Author LiuYong
// @Created at 2021-01-30
// @Modified at 2021-01-30
package setting

import (
	"common/tools/logging"
	"testing"
)

func TestSetting(t *testing.T) {
	logging.Setup()
	*configPath = "../../config/"
	Setup()
	t.Log(Consul)
}
