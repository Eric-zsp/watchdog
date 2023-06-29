package main

import (
	"os"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/jobs"
	gologs "github.com/cn-joyconn/gologs"
	"github.com/cn-joyconn/goutils/filetool"
	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	jobs.InitTask()
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
		gologs.GetLogger("default").Fatal(err.Error())
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			gologs.GetLogger("default").Info("服务安装成功" + global.AppConf.Name)
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			gologs.GetLogger("default").Info("服务卸载成功")
			return
		}
	}

	err = s.Run()
	if err != nil {
		gologs.GetLogger("default").Error(err.Error())
	}
}

func printLogo() {
	gologs.GetLogger("default").Info(edccLogo)
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
