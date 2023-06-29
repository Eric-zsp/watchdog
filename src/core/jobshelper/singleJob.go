package jobshelper

import (
	"sync"

	cron "github.com/Eric-zsp/cron/v3"
	gologs "github.com/cn-joyconn/gologs"
)

type SingleJobExcete interface {
	Execute(sj *SingleJob)
}
type SingleJob struct {
	ID      cron.EntryID
	JobKey  string
	JobName string
	Corn    string
	SJE     *SingleJobExcete
	runing  bool
	lock    sync.Mutex
}

func (sj *SingleJob) Run() {
	if !sj.runing {
		sj.lock.Lock()
		defer sj.lock.Unlock()
		if !sj.runing {
			sj.runing = true
			sj.doRun()
			sj.runing = false
		}

	}
}

func (sj *SingleJob) doRun() {
	defer func() {
		if err := recover(); err != nil {
			gologs.GetLogger("default").Error("roomcd:" + sj.JobName + ",执行采集任务异常，,错误信息：" + err.(error).Error())
		}
	}()
	(*sj.SJE).Execute(sj)
}

// func (sj *SingleJob) Execute() {
// 	fmt.Println("1111")
// }
