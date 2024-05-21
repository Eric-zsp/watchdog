package main

import (
	"flag"
	"os"
	"os/exec"
	"runtime"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/handle"
	"github.com/Eric-zsp/watchdog/src/core/jobs"
	"github.com/Eric-zsp/watchdog/src/core/utils/files"
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
		} else {

			// 标记
			var tags string
			// 代理行为
			var agent string
			// 服务名称
			var svrID string
			// 服务根目录
			var baseDir string
			// 服务下载包url
			var durl string
			// 文件、目录替换清单
			var listPath string

			flag.StringVar(&tags, "tags", "", "标记")
			flag.StringVar(&agent, "agent", "", "代理行为")
			flag.StringVar(&svrID, "svrID", "", "服务名称")
			flag.StringVar(&baseDir, "baseDir", "", "服务根目录")
			flag.StringVar(&durl, "url", "", "服务下载包url")
			flag.StringVar(&listPath, "listPath", "", "文件、目录替换清单")
			flag.Parse()

			gologs.GetLogger("default").Sugar().Info("启动参数：-tags=", tags, " -agent=", agent, " -svrID=", svrID, " -baseDir=", baseDir, " -durl=", durl, " -listPath=", listPath)
			if agent == "appRestart" {
				app_restart(svrID)
				return
			} else if agent == "appUpdate" {
				app_update(svrID, baseDir, durl, listPath)
				return
			}
		}
	}

	err = s.Run()
	if err != nil {
		gologs.GetLogger("default").Sugar().Error(err.Error())
	}
}

func app_restart(svrID string) {
	defer func() {
		if err := recover(); err != nil {
			gologs.GetLogger("default").Sugar().Error("the agent restart app happen error,the error msg:" + err.(error).Error())
		}
	}()
	gologs.GetLogger("default").Sugar().Info("the agent begin restart app")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// cmdStr := global.WindowsCMDAdminAuth + "net stop " + svrID + " & net start " + svrID
		cmdStr, e1 := files.SaveFileByes("restartApp.bat", []byte(global.WindowsCMDAdminAuth+"net stop "+svrID+" & net start "+svrID))
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Error("the agent  restarting app happen error,the error msg:" + e1.Error())
		}
		gologs.GetLogger("default").Sugar().Info(cmdStr)
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmdStr := "systemctl restart " + svrID
		gologs.GetLogger("default").Sugar().Info(cmdStr)
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}
	e2 := cmd.Run()
	if e2 != nil {
		gologs.GetLogger("default").Sugar().Error("the agent  restarting app happen error,the error msg:" + e2.Error())
	}
	gologs.GetLogger("default").Sugar().Info("the agent  restart app finished")
}

func app_update(svrID string, baseDir string, durl string, listPath string) {
	defer func() {
		if err := recover(); err != nil {
			gologs.GetLogger("default").Sugar().Error("the agent update app happen error,the error msg:" + err.(error).Error())
		}
	}()
	do := &handle.DogUpgradeOP{}
	do.DoUpgrade(svrID, baseDir, durl, listPath)
	gologs.GetLogger("default").Sugar().Info("the agent  update app finished")
}

func printLogo() {
	gologs.GetLogger("default").Sugar().Info(edccLogo)
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
