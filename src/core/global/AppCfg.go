package global

import (
	"sync"
	"time"

	gologs "github.com/cn-joyconn/gologs"
	"github.com/cn-joyconn/goutils/filetool"
	"github.com/json-iterator/go/extra"
	yaml "gopkg.in/yaml.v2"
)

var l sync.RWMutex
var appConfigPath string

var AppConf *AppConfig

var WaitExisit = false

var StartFinalStationTime time.Time

const (
	ConfigFinalconfiglogApiUrl      = "ConfigFinalconfiglogApi/"
	FinalUploadDataApiUrl           = "FinalUploadDataApi/"
	FinalUploadByCompressDataApiUrl = "FinalUploadByCompressDataApi/"
)

type Service struct {
	CheckName string `json:"checkName" yaml:"checkName"` //监测项名称
	CheckCorn string `json:"checkCorn" yaml:"checkCorn"` //监测周期
	CheckType int    `json:"checkType" yaml:"checkType"` //监测类型 1url(get) 2url(post)
	CheckAddr string `json:"checkAddr" yaml:"checkAddr"` //监测地址
	ErrOp     string `json:"errOp" yaml:"errOp"`         //监测失败时需要执行的动作(cmd/shell)
	ErrOpSpan int64  `json:"errOpSpan" yaml:"errOpSpan"` //监测失败时执行动作的最小时间间隔 单位秒
}

// 应用配置
type AppConfig struct {
	Name        string    `json:"name" yaml:"name"`               //应用名称
	DisplayName string    `json:"displayName" yaml:"displayName"` //服务显示名称
	Description string    `json:"description" yaml:"description"` //服务描述
	WebPort     int       `json:"webport" yaml:"webport"`         //web服务监听端口
	Services    []Service `json:"services" yaml:"services"`       //监测的服务
}

func init() {
	StartFinalStationTime = time.Now()
	extra.RegisterFuzzyDecoders()
}

// 加载程序配置
func InitAppConf(configPath string) {
	if filetool.IsExist(configPath) {
		configBytes, err := filetool.ReadFileToBytes(configPath)
		if err != nil {
			gologs.GetLogger("default").Error("读取" + configPath + "文件错误。" + err.Error())
			return
		}
		gologs.GetLogger("default").Info("成功读取" + configPath + "文件")
		// var appconfigs AppConfig
		_appconfigs := &AppConfig{
			WebPort: 7788,
		}
		err = yaml.Unmarshal(configBytes, _appconfigs)
		if err != nil {
			gologs.GetLogger("default").Error("解析" + configPath + "文件失败" + err.Error())
			return
		}
		gologs.GetLogger("default").Info("成功解析" + configPath + "文件")
		appConfigPath = configPath
		AppConf = _appconfigs
	} else {
		gologs.GetLogger("default").Error("未找到" + configPath)
		return
	}
}

// 保存程序配置
func SaveAppConfig() {
	appconfs := &AppConf
	saveYamlConfig(appconfs, appConfigPath)
}

// 保存yaml格式文件
func saveYamlConfig(in interface{}, configPath string) {
	l.Lock()
	defer l.Unlock()
	configBytes, err := yaml.Marshal(in)
	if err != nil {
		// fmt.Println(string(configBytes))
		gologs.GetLogger("default").Error(err.Error())
		return
	}
	_, err = filetool.WriteBytesToFile(configPath, configBytes)
	if err != nil {
		gologs.GetLogger("default").Error(err.Error())
		return
	}
}
