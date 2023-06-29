package jobs

import (
	"os/exec"
	"runtime"
	"time"

	global "github.com/Eric-zsp/watchdog/src/core/global"
	jobshelper "github.com/Eric-zsp/watchdog/src/core/jobsHelper"
	"github.com/Eric-zsp/watchdog/src/core/utils"
	gologs "github.com/cn-joyconn/gologs"
)

type ServiceCheckJob struct {
	Service  global.Service
	lastTime int64
}

func (handle *ServiceCheckJob) Execute(sj *jobshelper.SingleJob) {
	if global.WaitExisit {
		return
	}
	gologs.GetLogger("default").Debug("ServiceCheckJob 开始运行")
	if handle.Service.CheckType == 1 {
		_, result, err := utils.HttpGet(handle.Service.CheckAddr, "", nil, 10)
		if err != nil || result == "" {
			if (time.Now().Unix() - handle.lastTime) > handle.Service.ErrOpSpan {
				handle.lastTime = time.Now().Unix()
				var cmd *exec.Cmd
				if runtime.GOOS == "linux" {
					cmd = exec.Command("/bin/bash", "/c", handle.Service.ErrOp)
				} else {
					//windows 下
					cmd = exec.Command("cmd.exe", "/c", handle.Service.ErrOp)
				}
				err1 := cmd.Run()
				if err1 != nil {
					gologs.GetLogger("default").Error(handle.Service.CheckName + "  恢复服务发生错误：" + err1.Error())
				}
			}

		}
	}
	gologs.GetLogger("default").Debug("FinalStationConfigUpdateJob 运行结束")
}
