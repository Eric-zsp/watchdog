package main

import (
	"flag"
	"os"
	"os/exec"
	"runtime"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/handle"
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
		} else if os.Args[1] == "remove" {
			s.Uninstall()
			gologs.GetLogger("default").Info("服务卸载成功")
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
		gologs.GetLogger("default").Error(err.Error())
	}
}

func app_restart(svrID string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(global.UpgradeFileDir + "net stop " + svrID + " & net start " + svrID)
		cmd.Run()
		// cmd = exec.Command("copy", op.getCurrentCfgDir()+"*", dir)
	} else {
		cmd := exec.Command("systemctl restart " + svrID)
		cmd.Run()
	}
}

func app_update(svrID string, baseDir string, durl string, listPath string) {
	do := &handle.DogUpgradeOP{}
	do.DoUpgrade(svrID, baseDir, durl, listPath)
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
