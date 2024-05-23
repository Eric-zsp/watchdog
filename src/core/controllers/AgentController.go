package controllers

import (
	"os/exec"
	"runtime"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/handle"
	"github.com/Eric-zsp/watchdog/src/core/utils/files"
	gologs "github.com/cn-joyconn/gologs"
	"github.com/gin-gonic/gin"
)

type AgentParamController struct {
	BaseController
}

func (controller *AgentParamController) Agent_status(c *gin.Context) {

	controller.ApiSuccess(c, "the agent is live", 1)
}

func (controller *AgentParamController) App_restart(c *gin.Context) {
	svrID := c.PostForm("svrID")
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
	controller.ApiSuccess(c, "the agent  restart app finished", 1)
}

func (controller *AgentParamController) App_upgrade(c *gin.Context) {
	svrID := c.PostForm("svrID")
	baseDir := c.PostForm("baseDir")
	durl := c.PostForm("durl")
	listPath := c.PostForm("listPath")
	do := &handle.DogUpgradeOP{}
	do.DoUpgrade(svrID, baseDir, durl, listPath)
	gologs.GetLogger("default").Sugar().Info("the agent  upgrade app finished")
	controller.ApiSuccess(c, "the agent  upgrade app finished", 1)
}
