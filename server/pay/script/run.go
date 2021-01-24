// @Author LiuYong
// @Created at 2021-01-23
// @Modified at 2021-01-23
package script

import (
	"os/exec"
	"pay/tools/logging"
	"runtime"
)

func Setup() {
	systemType := runtime.GOOS
	var cmd *exec.Cmd
	switch systemType {
	case "linux":
		cmd = exec.Command("./script/linux-swagger.sh", "")
	case "windows":
		cmd = exec.Command(".\\script\\win-swagger.cmd", "")
	default:
		cmd = exec.Command(".\\script\\win-swagger.cmd", "")

	}
	err := cmd.Run()
	if err != nil {
		logging.Error("运行cmd错误", err)
	}
}
