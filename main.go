package main

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/initialize"
	gologs "github.com/cn-joyconn/gologs"
	"github.com/cn-joyconn/goutils/filetool"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {

	initialize.Init(func(e *gin.Engine) bool {
		timer1()
		// go testLog()
		return true
	})
	// 代码写在这儿
}

func (p *program) Stop(s service.Service) error {
	return nil
}

/**
* MAIN函数，程序入口
 */

func main() {
	printLogo()
	selfDir := filetool.SelfDir()
	appConfigPath := selfDir + "/conf/app.yml"
	global.InitAppConf(appConfigPath)

	svcConfig := &service.Config{
		Name:        global.AppConf.Name,        //服务显示名称
		DisplayName: global.AppConf.DisplayName, //服务名称
		Description: global.AppConf.Description, //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		gologs.GetLogger("default").Sugar().Fatal(err.Error())
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			gologs.GetLogger("default").Sugar().Info("服务安装成功" + global.AppConf.Name)
			return
		} else if os.Args[1] == "remove" {
			s.Uninstall()
			gologs.GetLogger("default").Sugar().Info("服务卸载成功")
			return
		}
	}

	err = s.Run()
	if err != nil {
		gologs.GetLogger("default").Sugar().Error(err.Error())
	}
}
func timer1() {
	timer1 := time.NewTicker(5 * time.Minute)
	// timer1 := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-timer1.C
			testTimer1()
		}
	}()

	// timer1 := time.NewTimer(5 * time.Second)
	// select {
	// case <-timer1.C:
	// 	testTimer1()
	// }

}
func testTimer1() {
	var m runtime.MemStats
	// var c runtime.
	runtime.ReadMemStats(&m)
	gologs.GetLogger("system").Sugar().Infof("%d Kb", m.Alloc/1024)

}
func printLogo() {
	logLine := strings.Split(edccLogo, "\n")
	for _, l := range logLine {
		gologs.GetLogger("default").Sugar().Info(l)
	}

	// gologs.GetLogger("[joyconn] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
	// 	utils.GlobalObject.Version,
	// 	utils.GlobalObject.MaxConn,
	// 	utils.GlobalObject.MaxPacketSize)
}

var edccLogo = `         
┌───────────────────────────────────────────────────────────────┐
      ██  ██████  ██    ██  ██████  ██████  ███    ██ ███    ██ 
      ██ ██    ██  ██  ██  ██      ██    ██ ████   ██ ████   ██ 
      ██ ██    ██   ████   ██      ██    ██ ██ ██  ██ ██ ██  ██ 
 ██   ██ ██    ██    ██    ██      ██    ██ ██  ██ ██ ██  ██ ██ 
  █████   ██████     ██     ██████  ██████  ██   ████ ██   ████ 
└───────────────────────────────────────────────────────────────┘ `
