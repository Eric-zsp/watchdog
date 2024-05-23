package global

const (
	SUCCESS     = 100001 //操作成功
	LoginSucess = 100002 //登录成功

	ERROR          = 100101 //"操作失败"
	ParamsError    = 100102 //参数错误
	NoResult       = 100103 //没有结果
	NotFound       = 100104 //没有结果
	ServiceError   = 100105 //服务器错误
	NoRule         = 100106 //没有权限
	CehckCodeError = 100107 //验证码错误
	VerifyFail     = 100108 //参数签名验证失败
	DbError        = 100109 //数据访问错误
	FiledRepeat    = 100110 //字段重复

	NoLogin        = 100201 //用户未登录
	TokenNotExist  = 100202 //令牌不存在
	TokenFail      = 100203 //令牌错误
	UserLocck      = 100204 //用户锁定
	LoginIdError   = 100205 //登录id错误
	LoginPassError = 100206 //密码错误
	LoginFail      = 100207 //登录失败
	UserNameExisit = 100208 //用户名已存在
	EmailExisit    = 100209 //邮箱已存在
	PhoneExisit    = 100210 //手机号已存在
	PwdWeek        = 100211 //密码强度不足

	ModbusErr = 200001 //通讯失败
)

// 响应参数结构体
type Response struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Data     interface{} `json:"data"`
	Url      string      `json:"url"`
	Wait     int         `json:"wait"`
	AllCount int64       `json:"allcount"`
}
