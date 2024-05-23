package utils

import (
	"strings"

	gin "github.com/gin-gonic/gin"
)

// PC（包括iE、google、sarfri）
var _PC = []string{"windows nT", "macintosh"}

// 安卓
var _Android = []string{"android"}

// IOS（"iPhone", "iPod", "iPad" ）
var _IOS = []string{"iphone", "ipod", "ipad"}

// windows phone
var _WP = []string{"windows phone"}

// 微信
var _WX = []string{"micromessenger", "wechat", "miniprogram"}

// 微信浏览器
var _WXBrowswr = []string{"micromessenger", "wechat"}

// 微信小程序
var _WXMiniApp = []string{"miniprogram"}

// 支付宝
var alipayclient = []string{"alipayclient"}

// ajax判断标准（http请求头中）
const ajaxRequestHeader = "XMLHttpRequest"

// GetClientType 获取客户端类型----pc还是移动端
// return 1:pc/未知    0:mobile
func GetClientType(ctx *gin.Context) int {
	userAgent := ctx.GetHeader("USER-AGENT")
	userAgent = strings.ToLower(userAgent)
	if len(userAgent) > 0 {
		pcList := make([]string, 0)
		pcList = append(pcList, _PC...)

		moblieList := make([]string, 0)
		moblieList = append(moblieList, _Android...)
		moblieList = append(moblieList, _IOS...)
		moblieList = append(moblieList, _WP...)
		moblieList = append(moblieList, _WX...)
		moblieList = append(moblieList, _WXBrowswr...)
		moblieList = append(moblieList, _WXMiniApp...)
		moblieList = append(moblieList, alipayclient...)
		//判断是否是pc端
		for _, source := range pcList {
			if strings.Index(userAgent, source) > -1 {
				return 1
			}
		}
		//判断是否是移动端。
		for _, source := range moblieList {
			if strings.Index(userAgent, source) > -1 {
				return 0
			}
		}

	}
	return 1
}

// IsAjax 判端是否为ajax请求
func IsAjax(ctx *gin.Context) bool {
	result := false
	xrequest := ctx.GetHeader("X-Requested-With")
	if ajaxRequestHeader == xrequest {
		result = true
	}
	if !result {
		contenType := ctx.ContentType()
		if strings.HasPrefix(contenType, "multipart/form-data") {
			result = true
		}
	}
	return result
}

// IsPage 判端是否为ajax请求
func IsPage(ctx *gin.Context) bool {
	result := false
	contenType := ctx.ContentType()
	if strings.HasPrefix(contenType, "text/html") {
		result = true
	}
	return result
}
