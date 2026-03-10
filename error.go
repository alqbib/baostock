package baostock

import "fmt"

// 错误代码
const (
	ErrSuccess              = "0"
	ErrNoLogin              = "10001001"
	ErrUsernameOrPassword   = "10001002"
	ErrGetUserInfoFail      = "10001003"
	ErrClientVersionExpire  = "10001004"
	ErrLoginCountLimit      = "10001005"
	ErrAccessInsufficience  = "10001006"
	ErrNeedActivate         = "10001007"
	ErrUsernameEmpty        = "10001008"
	ErrPasswordEmpty        = "10001009"
	ErrLogoutFail           = "10001010"
	ErrBlacklistUser        = "10001011"
	ErrSocketErr            = "10002001"
	ErrConnectFail          = "10002002"
	ErrConnectTimeout       = "10002003"
	ErrRecvConnectionClosed = "10002004"
	ErrSendSockFail         = "10002005"
	ErrSendSockTimeout      = "10002006"
	ErrRecvSockFail         = "10002007"
	ErrRecvSockTimeout      = "10002008"
	ErrParseDataErr         = "10004001"
	ErrUngzipDataFail       = "10004002"
	ErrUnknownErr           = "10004003"
	ErrOutOfBounds          = "10004004"
	ErrInparamEmpty         = "10004005"
	ErrParamErr             = "10004006"
	ErrStartDateErr         = "10004007"
	ErrEndDateErr           = "10004008"
	ErrStartBigthanEnd      = "10004009"
	ErrDateErr              = "10004010"
	ErrCodeInvalided        = "10004011"
	ErrIndicatorInvalided   = "10004012"
	ErrBeyondDateSupport    = "10004013"
	ErrMixedCodesMarket     = "10004014"
	ErrNoSupportCodesMarket = "10004015"
	ErrOrderToUpperLimit    = "10004016"
	ErrNoSupportOrderInfo   = "10004017"
	ErrIndicatorRepeat      = "10004018"
	ErrMessageError         = "10004019"
	ErrMessageCodeError     = "10004020"
	ErrSystemError          = "10005001"
)

// 错误信息映射
var errorMessages = map[string]string{
	ErrSuccess:              "成功",
	ErrNoLogin:              "用户未登录",
	ErrUsernameOrPassword:   "用户名或密码错误",
	ErrGetUserInfoFail:      "获取用户信息失败",
	ErrClientVersionExpire:  "客户端版本号过期",
	ErrLoginCountLimit:      "账号登录数达到上限",
	ErrAccessInsufficience:  "用户权限不足",
	ErrNeedActivate:         "需要登录激活",
	ErrUsernameEmpty:        "用户名为空",
	ErrPasswordEmpty:        "密码为空",
	ErrLogoutFail:           "用户登出失败",
	ErrBlacklistUser:        "黑名单用户",
	ErrSocketErr:            "网络错误",
	ErrConnectFail:          "网络连接失败",
	ErrConnectTimeout:       "网络连接超时",
	ErrRecvConnectionClosed: "网络接收时连接断开",
	ErrSendSockFail:         "网络发送失败",
	ErrSendSockTimeout:      "网络发送超时",
	ErrRecvSockFail:         "网络接收错误",
	ErrRecvSockTimeout:      "网络接收超时",
	ErrParseDataErr:         "解析数据错误",
	ErrUngzipDataFail:       "gzip解压失败",
	ErrUnknownErr:           "客户端未知错误",
	ErrOutOfBounds:          "数组越界",
	ErrInparamEmpty:         "传入参数为空",
	ErrParamErr:             "参数错误",
	ErrStartDateErr:         "起始日期格式不正确",
	ErrEndDateErr:           "截止日期格式不正确",
	ErrStartBigthanEnd:      "起始日期大于终止日期",
	ErrDateErr:              "日期格式不正确",
	ErrCodeInvalided:        "无效的证券代码",
	ErrIndicatorInvalided:   "无效的指标",
	ErrBeyondDateSupport:    "超出日期支持范围",
	ErrMixedCodesMarket:     "不支持的混合证券品种",
	ErrNoSupportCodesMarket: "不支持的证券代码品种",
	ErrOrderToUpperLimit:    "交易条数超过上限",
	ErrNoSupportOrderInfo:   "不支持的交易信息",
	ErrIndicatorRepeat:      "指标重复",
	ErrMessageError:         "消息格式不正确",
	ErrMessageCodeError:     "错误的消息类型",
	ErrSystemError:          "系统级别错误",
}

// Error 表示 BaoStock 错误,可以使用类型断言检查特定错误: errors.As(err, &baostock.Error{})
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	if msg, ok := errorMessages[e.Code]; ok {
		return fmt.Sprintf("%s: %s", e.Code, msg)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
