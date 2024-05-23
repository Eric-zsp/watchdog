package controllers

import (
	"net/http"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/gin-gonic/gin"
	// gologs "github.com/cn-joyconn/gologs"
	// "strconv"
)

type BaseController struct {
	//gin.Context
}

// JSON输出
func (bc *BaseController) ApiJson(c *gin.Context, code int, msg string, data interface{}, allcount int64) {
	if data == nil {
		data = ""
	}
	// c.ServeJSON()
	// c.StopRun()
	c.JSON(http.StatusOK, &global.Response{
		Code:     code,
		Msg:      msg,
		Data:     data,
		Url:      "",
		Wait:     0,
		AllCount: allcount,
	})
}

// 返回成功的API成功
func (bc *BaseController) ApiSuccess(c *gin.Context, msg string, data interface{}) {
	bc.ApiJson(c, global.SUCCESS, msg, data, 0)
}

// 返回成功的API成功
func (bc *BaseController) ApiDataList(c *gin.Context, msg string, data interface{}, allcount int64) {
	bc.ApiJson(c, global.SUCCESS, msg, data, allcount)
}

// 返回失败的API请求
func (bc *BaseController) ApiError(c *gin.Context, msg string, data interface{}) {
	bc.ApiJson(c, global.ERROR, msg, data, 0)
}

// 返回失败且带code的API请求
func (bc *BaseController) ApiErrorCode(c *gin.Context, msg string, data interface{}, code int) {
	bc.ApiJson(c, code, msg, data, 0)
}
