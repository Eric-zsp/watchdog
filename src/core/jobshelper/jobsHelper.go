package jobsHelper

import (
	"errors"
	"strings"
	"sync"

	cron "github.com/Eric-zsp/cron/v3"
)

var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.DowOptional | cron.Descriptor)
var cornInstance *cron.Cron
var singleJobInstances map[string]*SingleJob //单例字典
var singleJobInstancesLock sync.Mutex        // 锁对象

func init() {
	singleJobInstances = make(map[string]*SingleJob, 0)
	cornInstance = cron.New(cron.WithParser(secondParser), cron.WithChain())
	StartJobs()
}

/**
 * @Description: 添加一个定时任务
 *
 */
func AddJob(JobKey string, JobName string, Corn string, sje *SingleJobExcete) (bool, error) {
	singleJobInstancesLock.Lock()
	defer singleJobInstancesLock.Unlock()
	_, ok := singleJobInstances[JobKey]
	if ok {
		return false, errors.New("job's JobKey has exisit")
	} else {
		cronStr := getCron(Corn)
		sj := &SingleJob{
			JobKey:  JobKey,
			JobName: JobName,
			Corn:    cronStr,
			SJE:     sje,
		}
		ID, err := cornInstance.AddJob(sj.Corn, sj)
		if err != nil {
			return false, err
		}
		sj.ID = ID
		singleJobInstances[sj.JobKey] = sj
		return true, nil
	}
}

/**
 * @Description: 修改一个任务的触发时间
 *
 *  JobKey job实例键
 *  cron   时间设置，参考quartz说明文档
 *               @return 1：修改成功，0：修改失败，-1:任务不存在
 */
func ModifyJobTime(JobKey string, cronStr string) (bool, error) {
	singleJobInstancesLock.Lock()
	defer singleJobInstancesLock.Unlock()
	sj, ok := singleJobInstances[JobKey]
	if !ok {
		return false, errors.New("job's JobKey not exisit")
	} else {
		// ent := cornInstance.Entry(sj.ID)
		cronStr = getCron(cronStr)
		cornInstance.Remove(sj.ID)
		ID, err := cornInstance.AddJob(cronStr, sj)
		if err != nil {
			return false, err
		}
		sj.ID = ID
		sj.Corn = cronStr
		// cornInstance.start
		return true, nil
	}
}

/**
 * @Description: 移除一个任务
 *
 *  JobKey
 */
func RemoveJob(JobKey string) {
	singleJobInstancesLock.Lock()
	defer singleJobInstancesLock.Unlock()
	sj, ok := singleJobInstances[JobKey]
	if ok {
		cornInstance.Remove(sj.ID)
		delete(singleJobInstances, JobKey)
	}
}

/**
 * @title:启动所有定时任务
 */
func StartJobs() {
	cornInstance.Start()
}

/**
 * @title:关闭所有定时任务
 */
func ShutdownJobs() {
	cornInstance.Stop()
}

func getCron(cron string) string {
	if strings.HasPrefix(cron, "@every ") {
		return cron
	} else if strings.HasSuffix(cron, "ms") {
		return "@every " + cron
	} else if strings.HasSuffix(cron, "s") {
		return "@every " + cron
	} else if strings.HasSuffix(cron, "m") {
		return "@every " + cron
	} else if strings.HasSuffix(cron, "h") {
		return "@every " + cron
	} else if strings.HasSuffix(cron, "d") {
		return "@every " + cron
	} else {
		cs := strings.Split(cron, " ")
		if len(cs) > 6 {
			cron = strings.Join(cs[:6], " ")
		}
	}
	return cron
}
