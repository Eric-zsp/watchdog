package jobs

import (
	global "github.com/Eric-zsp/watchdog/src/core/global"
	jobshelper "github.com/Eric-zsp/watchdog/src/core/jobsHelper"
	gologs "github.com/cn-joyconn/gologs"
	"github.com/cn-joyconn/goutils/strtool"
)

func InitTask() {
	if global.AppConf.Services != nil {
		for _, serviceItem := range global.AppConf.Services {
			if !strtool.IsBlank(serviceItem.CheckCorn) {
				var sjt1 jobshelper.SingleJobExcete
				sjt1 = &ServiceCheckJob{
					Service: serviceItem,
				}
				added, err := jobshelper.AddJob(serviceItem.CheckName, serviceItem.CheckName, serviceItem.CheckCorn, &sjt1)
				if err != nil {
					gologs.GetLogger("default").Debug("Check service Job 错误，" + err.Error())
				} else if !added {
					gologs.GetLogger("default").Debug("Check service Job 错误")
				}
			}
		}
	}
}
