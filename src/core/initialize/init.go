package initialize

import (
	"strconv"

	"github.com/Eric-zsp/watchdog/src/core/controllers"
	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/jobs"
	"github.com/Eric-zsp/watchdog/src/core/middleware"
	gologs "github.com/cn-joyconn/gologs"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go/extra"
)

type ExtInit func(*gin.Engine) bool

func init() {
	extra.RegisterFuzzyDecoders()
}

// Init 初始化
func Init(f ExtInit) {

	jobs.InitTask()
	Router := InitServer()
	RegistorRouters(Router)

	if f(Router) {
		Router.Run(":" + strconv.Itoa(global.AppConf.WebPort))
		// Router.RunTLS(":"+strconv.Itoa(global.AppConf.WebPort), selfDir+"/conf/ssl/ssl.pem", selfDir+"/conf/ssl/ssl.key")
	}

}

func InitServer() *gin.Engine {
	// if !AppConf.Debug {
	// 	gin.SetMode(gin.ReleaseMode)
	// 	gin.DefaultWriter = ioutil.Discard
	// }
	// selfDir := filetool.SelfDir()
	var Router = gin.Default()

	// 日志
	Router.Use(middleware.GinLogger())
	gologs.GetLogger("default").Info("use middleware Logger")

	// https
	// Router.Use(middleware.LoadTls()) // 打开就能玩https了

	//Error
	Router.NoMethod(middleware.HandleNotFound)
	Router.NoRoute(middleware.HandleNotFound)
	Router.Use(middleware.ErrHandler())
	gologs.GetLogger("default").Info("use middleware ErrorHandle")

	return Router

}

// RegistorRouters 初始化总路由
func RegistorRouters(Router *gin.Engine) {

	contextRouter := Router.Group("/")
	// // 方便统一添加路由组前缀 多服务器上线使用
	publicGroup := contextRouter.Group("")
	agentController := &controllers.AgentParamController{}
	publicGroup.POST("/api/app/agent/status", agentController.Agent_status)
	publicGroup.POST("/api/app/agent/restart", agentController.App_restart)
	publicGroup.POST("/api/app/agent/upgrade", agentController.App_upgrade)

	gologs.GetLogger("default").Info("router register success")
	gologs.GetLogger("system").Info("router register success")
}
