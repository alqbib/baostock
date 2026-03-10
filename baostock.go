// Package baostock 提供 BaoStock 证券数据平台的 Go 语言客户端。
// BaoStock 是一个免费、开源的证券数据平台，提供完整的A股历史数据、实时行情、财务数据等。
//
// 使用示例:
//
//	client := baostock.NewClient()
//	if err := client.Login(context.Background()); err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Logout()
//
//	data, err := client.QueryHistoryKDataPlus(context.Background(),
//	    &baostock.HistoryKDataRequest{
//	        Code:        "sh.600000",
//	        Fields:      "date,code,open,high,low,close,volume",
//	        StartDate:   "2023-01-01",
//	        EndDate:     "2023-12-31",
//	        Frequency:   baostock.FrequencyDaily,
//	        AdjustFlag:  baostock.AdjustFlagNoAdjust,
//	    })
package baostock

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// 协议常量
const (
	// 服务器配置
	DefaultServerHost = "www.baostock.com"
	DefaultServerPort = 10030
	ClientVersion     = "00.8.90"

	// 消息分隔符
	MessageSplit = "\x01" // ASCII 0x01
	Delimiter    = "\n"   // ASCII 0x0A

	// 消息结构
	MessageHeaderLength      = 21
	MessageHeaderBodyLength  = 10

	// 分页
	DefaultPerPageCount = 10000

	// 证券代码
	StockCodeLength = 9

	// 默认值
	DefaultStartDate = "2015-01-01"
)

// 消息类型代码
const (
	// 登录/登出
	MsgTypeLoginRequest  = "00"
	MsgTypeLoginResponse = "01"
	MsgTypeLogoutRequest = "02"
	MsgTypeLogoutResponse = "03"
	MsgTypeError         = "04"

	// K线数据
	MsgTypeGetKDataRequest      = "11"
	MsgTypeGetKDataResponse     = "12"
	MsgTypeGetKDataPlusRequest  = "95"
	MsgTypeGetKDataPlusResponse = "96"

	// 财务数据
	MsgTypeQueryDividendDataRequest    = "13"
	MsgTypeQueryDividendDataResponse   = "14"
	MsgTypeAdjustFactorRequest         = "15"
	MsgTypeAdjustFactorResponse        = "16"
	MsgTypeProfitDataRequest           = "17"
	MsgTypeProfitDataResponse          = "18"
	MsgTypeOperationDataRequest        = "19"
	MsgTypeOperationDataResponse       = "20"
	MsgTypeQueryGrowthDataRequest      = "21"
	MsgTypeQueryGrowthDataResponse     = "22"
	MsgTypeQueryDupontDataRequest      = "23"
	MsgTypeQueryDupontDataResponse     = "24"
	MsgTypeQueryBalanceDataRequest     = "25"
	MsgTypeQueryBalanceDataResponse    = "26"
	MsgTypeQueryCashFlowDataRequest    = "27"
	MsgTypeQueryCashFlowDataResponse   = "28"

	// 公司公告
	MsgTypeQueryPerformanceExpressReportRequest  = "29"
	MsgTypeQueryPerformanceExpressReportResponse = "30"
	MsgTypeQueryForecastReportRequest            = "31"
	MsgTypeQueryForecastReportResponse           = "32"

	// 元数据查询
	MsgTypeQueryTradeDatesRequest  = "33"
	MsgTypeQueryTradeDatesResponse = "34"
	MsgTypeQueryAllStockRequest    = "35"
	MsgTypeQueryAllStockResponse   = "36"
	MsgTypeQueryStockBasicRequest  = "45"
	MsgTypeQueryStockBasicResponse = "46"

	// 板块信息
	MsgTypeQueryStockIndustryRequest  = "59"
	MsgTypeQueryStockIndustryResponse = "60"
	MsgTypeQueryStockConceptRequest   = "81"
	MsgTypeQueryStockConceptResponse  = "82"
	MsgTypeQueryStockAreaRequest      = "83"
	MsgTypeQueryStockAreaResponse     = "84"

	// 指数成分股
	MsgTypeQueryHS300StocksRequest  = "61"
	MsgTypeQueryHS300StocksResponse = "62"
	MsgTypeQuerySZ50StocksRequest   = "63"
	MsgTypeQuerySZ50StocksResponse  = "64"
	MsgTypeQueryZZ500StocksRequest  = "65"
	MsgTypeQueryZZ500StocksResponse = "66"

	// 特殊股票
	MsgTypeQueryTerminatedStocksRequest  = "67"
	MsgTypeQueryTerminatedStocksResponse = "68"
	MsgTypeQuerySuspendedStocksRequest   = "69"
	MsgTypeQuerySuspendedStocksResponse  = "70"
	MsgTypeQuerySTStocksRequest          = "71"
	MsgTypeQuerySTStocksResponse         = "72"
	MsgTypeQueryStarSTStocksRequest      = "73"
	MsgTypeQueryStarSTStocksResponse     = "74"

	// 市场板块
	MsgTypeQueryAMEStocksRequest  = "85"
	MsgTypeQueryAMEStocksResponse = "86"
	MsgTypeQueryGEMStocksRequest  = "87"
	MsgTypeQueryGEMStocksResponse = "88"

	// 市场互联互通
	MsgTypeQuerySHHKStocksRequest  = "89"
	MsgTypeQuerySHHKStocksResponse = "90"
	MsgTypeQuerySZHKStocksRequest  = "91"
	MsgTypeQuerySZHKStocksResponse = "92"

	// 风险警示
	MsgTypeQueryStockInRiskRequest  = "93"
	MsgTypeQueryStockInRiskResponse = "94"

	// 宏观经济数据
	MsgTypeQueryDepositRateDataRequest         = "47"
	MsgTypeQueryDepositRateDataResponse        = "48"
	MsgTypeQueryLoanRateDataRequest            = "49"
	MsgTypeQueryLoanRateDataResponse           = "50"
	MsgTypeQueryRequiredReserveRatioDataRequest = "51"
	MsgTypeQueryRequiredReserveRatioDataResponse = "52"
	MsgTypeQueryMoneySupplyDataMonthRequest    = "53"
	MsgTypeQueryMoneySupplyDataMonthResponse   = "54"
	MsgTypeQueryMoneySupplyDataYearRequest     = "55"
	MsgTypeQueryMoneySupplyDataYearResponse    = "56"
	MsgTypeQuerySHIBORDataRequest              = "57"
	MsgTypeQuerySHIBORDataResponse             = "58"
	MsgTypeQueryCPIDataRequest                 = "75"
	MsgTypeQueryCPIDataResponse                = "76"
	MsgTypeQueryPPIDataRequest                 = "77"
	MsgTypeQueryPPIDataResponse                = "78"
	MsgTypeQueryPMIDataRequest                 = "79"
	MsgTypeQueryPMIDataResponse                = "80"

	// 实时行情
	MsgTypeLoginRealTimeRequest    = "37"
	MsgTypeLoginRealTimeResponse   = "38"
	MsgTypeLogoutRealTimeRequest   = "39"
	MsgTypeLogoutRealTimeResponse  = "40"
	MsgTypeSubscriptionsRequest    = "41"
	MsgTypeSubscriptionsResponse   = "42"
	MsgTypeCancelSubscribeRequest  = "43"
	MsgTypeCancelSubscribeResponse = "44"
)

// 错误代码
const (
	ErrSuccess                = "0"
	ErrNoLogin               = "10001001"
	ErrUsernameOrPassword    = "10001002"
	ErrGetUserInfoFail       = "10001003"
	ErrClientVersionExpire   = "10001004"
	ErrLoginCountLimit       = "10001005"
	ErrAccessInsufficience   = "10001006"
	ErrNeedActivate          = "10001007"
	ErrUsernameEmpty         = "10001008"
	ErrPasswordEmpty         = "10001009"
	ErrLogoutFail            = "10001010"
	ErrBlacklistUser         = "10001011"
	ErrSocketErr             = "10002001"
	ErrConnectFail           = "10002002"
	ErrConnectTimeout        = "10002003"
	ErrRecvConnectionClosed  = "10002004"
	ErrSendSockFail          = "10002005"
	ErrSendSockTimeout       = "10002006"
	ErrRecvSockFail          = "10002007"
	ErrRecvSockTimeout       = "10002008"
	ErrParseDataErr          = "10004001"
	ErrUngzipDataFail        = "10004002"
	ErrUnknownErr            = "10004003"
	ErrOutOfBounds           = "10004004"
	ErrInparamEmpty          = "10004005"
	ErrParamErr              = "10004006"
	ErrStartDateErr          = "10004007"
	ErrEndDateErr            = "10004008"
	ErrStartBigthanEnd       = "10004009"
	ErrDateErr               = "10004010"
	ErrCodeInvalided         = "10004011"
	ErrIndicatorInvalided    = "10004012"
	ErrBeyondDateSupport     = "10004013"
	ErrMixedCodesMarket      = "10004014"
	ErrNoSupportCodesMarket  = "10004015"
	ErrOrderToUpperLimit     = "10004016"
	ErrNoSupportOrderInfo    = "10004017"
	ErrIndicatorRepeat       = "10004018"
	ErrMessageError          = "10004019"
	ErrMessageCodeError      = "10004020"
	ErrSystemError           = "10005001"
)

// 错误信息映射
var errorMessages = map[string]string{
	ErrSuccess:                "成功",
	ErrNoLogin:               "用户未登录",
	ErrUsernameOrPassword:    "用户名或密码错误",
	ErrGetUserInfoFail:       "获取用户信息失败",
	ErrClientVersionExpire:   "客户端版本号过期",
	ErrLoginCountLimit:       "账号登录数达到上限",
	ErrAccessInsufficience:   "用户权限不足",
	ErrNeedActivate:          "需要登录激活",
	ErrUsernameEmpty:         "用户名为空",
	ErrPasswordEmpty:         "密码为空",
	ErrLogoutFail:            "用户登出失败",
	ErrBlacklistUser:         "黑名单用户",
	ErrSocketErr:             "网络错误",
	ErrConnectFail:           "网络连接失败",
	ErrConnectTimeout:        "网络连接超时",
	ErrRecvConnectionClosed:  "网络接收时连接断开",
	ErrSendSockFail:          "网络发送失败",
	ErrSendSockTimeout:       "网络发送超时",
	ErrRecvSockFail:          "网络接收错误",
	ErrRecvSockTimeout:       "网络接收超时",
	ErrParseDataErr:          "解析数据错误",
	ErrUngzipDataFail:        "gzip解压失败",
	ErrUnknownErr:            "客户端未知错误",
	ErrOutOfBounds:           "数组越界",
	ErrInparamEmpty:          "传入参数为空",
	ErrParamErr:              "参数错误",
	ErrStartDateErr:          "起始日期格式不正确",
	ErrEndDateErr:            "截止日期格式不正确",
	ErrStartBigthanEnd:       "起始日期大于终止日期",
	ErrDateErr:               "日期格式不正确",
	ErrCodeInvalided:         "无效的证券代码",
	ErrIndicatorInvalided:    "无效的指标",
	ErrBeyondDateSupport:     "超出日期支持范围",
	ErrMixedCodesMarket:      "不支持的混合证券品种",
	ErrNoSupportCodesMarket:  "不支持的证券代码品种",
	ErrOrderToUpperLimit:     "交易条数超过上限",
	ErrNoSupportOrderInfo:    "不支持的交易信息",
	ErrIndicatorRepeat:       "指标重复",
	ErrMessageError:          "消息格式不正确",
	ErrMessageCodeError:      "错误的消息类型",
	ErrSystemError:           "系统级别错误",
}

// Frequency 表示K线频率
type Frequency string

const (
	Frequency5Min  Frequency = "5"  // 5分钟
	Frequency15Min Frequency = "15" // 15分钟
	Frequency30Min Frequency = "30" // 30分钟
	Frequency60Min Frequency = "60" // 60分钟
	FrequencyDaily Frequency = "d"  // 日线
	FrequencyWeek  Frequency = "w"  // 周线
	FrequencyMonth Frequency = "m"  // 月线
)

// AdjustFlag 表示复权类型
type AdjustFlag string

const (
	AdjustFlagBackward  AdjustFlag = "1" // 后复权
	AdjustFlagForward   AdjustFlag = "2" // 前复权
	AdjustFlagNoAdjust  AdjustFlag = "3" // 不复权
)

// Config 表示客户端配置
type Config struct {
	Host     string        // 服务器地址
	Port     int           // 服务器端口
	Username string        // 用户名
	Password string        // 密码
	Timeout  time.Duration // 超时时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:     DefaultServerHost,
		Port:     DefaultServerPort,
		Username: "anonymous",
		Password: "123456",
		Timeout:  30 * time.Second,
	}
}

// Client 表示 BaoStock 客户端
type Client struct {
	config    *Config
	conn      net.Conn
	reader    *bufio.Reader
	connected bool
	loggedIn  bool
	userID    string
}

// NewClient 使用默认配置创建新的 BaoStock 客户端
func NewClient() *Client {
	return NewClientWithConfig(DefaultConfig())
}

// NewClientWithConfig 使用自定义配置创建新的 BaoStock 客户端
func NewClientWithConfig(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}
	return &Client{
		config: config,
	}
}

// Connect 建立与服务器的连接
func (c *Client) Connect(ctx context.Context) error {
	if c.conn != nil {
		return nil // 已连接
	}

	dialer := &net.Dialer{
		Timeout: c.config.Timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", c.config.Host, c.config.Port))
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReaderSize(conn, 8192)
	c.connected = true
	return nil
}

// Disconnect 关闭与服务器的连接
func (c *Client) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	c.reader = nil
	c.connected = false
	c.loggedIn = false
	return err
}

// Login 与 BaoStock 服务器进行身份验证
func (c *Client) Login(ctx context.Context) error {
	if err := c.Connect(ctx); err != nil {
		return err
	}

	msgBody := fmt.Sprintf("login%s%s%s%s%s0", MessageSplit, c.config.Username, MessageSplit, c.config.Password, MessageSplit)

	resp, err := c.sendMessage(ctx, MsgTypeLoginRequest, msgBody)
	if err != nil {
		return err
	}

	result := &LoginResponse{}
	if err := parseLoginResponse(resp, result); err != nil {
		return err
	}

	if result.ErrorCode != ErrSuccess {
		return &Error{Code: result.ErrorCode, Message: result.ErrorMsg}
	}

	c.userID = result.UserID
	c.loggedIn = true
	return nil
}

// Logout 终止会话
func (c *Client) Logout(ctx context.Context) error {
	if !c.loggedIn {
		return nil
	}

	timestamp := time.Now().Format("20060102150405")
	msgBody := fmt.Sprintf("logout%s%s%s%s", MessageSplit, c.userID, MessageSplit, timestamp)

	_, err := c.sendMessage(ctx, MsgTypeLogoutRequest, msgBody)
	if err != nil {
		return err
	}

	c.loggedIn = false
	return c.Disconnect()
}

// sendMessage 向服务器发送消息并返回响应
func (c *Client) sendMessage(ctx context.Context, msgType, msgBody string) (*Response, error) {
	if !c.connected {
		return nil, errors.New("未连接到服务器")
	}

	// 构建消息头
	header := fmt.Sprintf("%s%s%s%s%s", ClientVersion, MessageSplit, msgType, MessageSplit, formatLength(len(msgBody)))

	// 计算CRC32
	fullMsg := header + msgBody
	crc32 := calculateCRC32(fullMsg)

	// 构建完整请求
	request := fmt.Sprintf("%s%s%d%s", fullMsg, MessageSplit, crc32, Delimiter)

	// 发送请求
	if _, err := c.conn.Write([]byte(request)); err != nil {
		return nil, fmt.Errorf("发送消息失败: %w", err)
	}

	// 接收响应
	response, err := c.receiveResponse()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// receiveResponse 接收并解析来自服务器的响应
func (c *Client) receiveResponse() (*Response, error) {
	var buffer bytes.Buffer

	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		buffer.WriteString(line)

		// 检查消息是否结束
		str := buffer.String()
		// 对于压缩响应，必须等待 <![CDATA[]]>\n 标记
		// 对于非压缩响应，只需要 \n
		// 先检查是否是压缩响应类型（通过检查消息头）
		if len(str) >= MessageHeaderLength {
			header := str[:MessageHeaderLength]
			headerParts := strings.Split(header, MessageSplit)
			if len(headerParts) >= 2 {
				msgType := headerParts[1]
				// 检查是否是压缩消息类型
				isCompressed := (msgType == MsgTypeGetKDataPlusResponse)

				if isCompressed && strings.HasSuffix(str, "<![CDATA[]]>\n") {
					break
				} else if !isCompressed && strings.HasSuffix(str, "\n") {
					break
				}
			}
		}
	}

	responseStr := buffer.String()

	// 解析消息头
	if len(responseStr) < MessageHeaderLength {
		return nil, errors.New("无效的响应: 消息过短")
	}

	header := responseStr[:MessageHeaderLength]
	body := responseStr[MessageHeaderLength:]

	// 检查是否压缩
	headerParts := strings.Split(header, MessageSplit)
	if len(headerParts) < 3 {
		return nil, errors.New("无效的响应头")
	}

	msgType := headerParts[1]

	// 处理压缩响应
	if msgType == MsgTypeGetKDataPlusResponse {
		// 对于压缩响应，数据格式为: header + [压缩数据]\x1[CRC32]<![CDATA[]]>\n
		// 需要从完整响应中找到 <![CDATA[ 的位置
		before, _, ok := strings.Cut(responseStr, "<![CDATA[")
		if !ok {
			// 未找到 <![CDATA[ 标记，尝试找到数据结束位置
			// 响应格式: header + compressed_data + \x01 + CRC32 + \n
			// 从后往前找第一个 \x01 作为CRC32分隔符
			lastNewline := strings.LastIndex(responseStr, Delimiter)

			if lastNewline > MessageHeaderLength {
				// 在换行符前找最后一个 \x01
				lastSplit := strings.LastIndex(responseStr[:lastNewline], MessageSplit)
				if lastSplit > MessageHeaderLength {
					compressedData := []byte(responseStr[MessageHeaderLength:lastSplit])

					decompressed, err := decompressData(compressedData)
					if err != nil {
						return nil, fmt.Errorf("解压失败: %w", err)
					}

					body = string(decompressed)
				} else {
					return nil, errors.New("压缩响应格式错误: 无法确定数据范围")
				}
			} else {
				return nil, errors.New("压缩响应格式错误: 响应过短")
			}
		} else {
			// 找到 <![CDATA[ 前面的 \x01 分隔符（CRC32前的分隔符）
			cdataEnd := strings.LastIndex(before, MessageSplit)
			if cdataEnd == -1 {
				return nil, errors.New("压缩响应格式错误: 未找到CRC32分隔符")
			}

			compressedData := []byte(responseStr[MessageHeaderLength:cdataEnd])

			decompressed, err := decompressData(compressedData)
			if err != nil {
				return nil, fmt.Errorf("解压失败: %w", err)
			}

			body = string(decompressed)
		}
	}

	return &Response{
		Header:     header,
		Body:       body,
		FullString: responseStr,
	}, nil
}

// Response 表示服务器响应
type Response struct {
	Header     string // 消息头
	Body       string // 消息体
	FullString string // 完整响应字符串
}

// LoginResponse 表示登录响应
type LoginResponse struct {
	ErrorCode string
	ErrorMsg  string
	Method    string
	UserID    string
}

// HistoryKDataRequest 表示历史K线数据请求
type HistoryKDataRequest struct {
	Code       string   // 证券代码
	Fields     string   // 字段列表
	StartDate  string   // 开始日期
	EndDate    string   // 结束日期
	Frequency  Frequency // K线频率
	AdjustFlag AdjustFlag // 复权标志
}

// HistoryKDataResponse 表示历史K线数据响应
type HistoryKDataResponse struct {
	ErrorCode     string   // 错误代码
	ErrorMsg      string   // 错误信息
	Method        string   // 方法名
	UserID        string   // 用户ID
	CurPageNum    string   // 当前页码
	PerPageCount  string   // 每页条数
	Data          [][]string // 数据
	Code          string   // 证券代码
	Fields        []string // 字段列表
	StartDate     string   // 开始日期
	EndDate       string   // 结束日期
	Frequency     string   // K线频率
	AdjustFlag    string   // 复权标志
}

// QueryHistoryKDataPlus 查询历史K线数据
func (c *Client) QueryHistoryKDataPlus(ctx context.Context, req *HistoryKDataRequest) (*HistoryKDataResponse, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if err := validateStockCode(req.Code); err != nil {
		return nil, err
	}

	startDate := req.StartDate
	if startDate == "" {
		startDate = DefaultStartDate
	}

	endDate := req.EndDate
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 规范化股票代码
	code := normalizeStockCode(req.Code)

	// query_history_k_data_plus\x1用户ID\x1页码\x1每页条数\x1证券代码\x1字段列表\x1开始日期\x1结束日期\x1频率\x1复权标志
	msgBody := fmt.Sprintf("query_history_k_data_plus%s%s%s1%s%d%s%s%s%s%s%s%s%s%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, code, MessageSplit,
		req.Fields, MessageSplit, startDate, MessageSplit,
		endDate, MessageSplit, req.Frequency, MessageSplit, req.AdjustFlag)

	resp, err := c.sendMessage(ctx, MsgTypeGetKDataPlusRequest, msgBody)
	if err != nil {
		return nil, err
	}

	return parseHistoryKDataResponse(resp)
}

// QueryTradeDates 查询交易日
func (c *Client) QueryTradeDates(ctx context.Context, startDate, endDate string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if startDate == "" {
		startDate = DefaultStartDate
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	msgBody := fmt.Sprintf("query_trade_dates%s%s%s1%s%d%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, startDate, MessageSplit, endDate)

	resp, err := c.sendMessage(ctx, MsgTypeQueryTradeDatesRequest, msgBody)
	if err != nil {
		return nil, err
	}

	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 7 {
		return nil, errors.New("无效的响应")
	}

	errorCode := bodyParts[0]
	if errorCode != ErrSuccess {
		errorMsg := ""
		if len(bodyParts) > 1 {
			errorMsg = bodyParts[1]
		}
		return nil, &Error{Code: errorCode, Message: errorMsg}
	}

	// 解析JSON数据
	var result struct {
		Data [][]string `json:"record"`
	}
	if err := json.Unmarshal([]byte(bodyParts[6]), &result); err != nil {
		return nil, fmt.Errorf("解析数据JSON失败: %w", err)
	}

	return result.Data, nil
}

// QueryAllStock 查询指定日期的所有股票
func (c *Client) QueryAllStock(ctx context.Context, date string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	msgBody := fmt.Sprintf("query_all_stock%s%s%s1%s%d%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, date)

	resp, err := c.sendMessage(ctx, MsgTypeQueryAllStockRequest, msgBody)
	if err != nil {
		return nil, err
	}

	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 7 {
		return nil, errors.New("无效的响应")
	}

	errorCode := bodyParts[0]
	if errorCode != ErrSuccess {
		errorMsg := ""
		if len(bodyParts) > 1 {
			errorMsg = bodyParts[1]
		}
		return nil, &Error{Code: errorCode, Message: errorMsg}
	}

	// 解析JSON数据
	var result struct {
		Data [][]string `json:"record"`
	}
	if err := json.Unmarshal([]byte(bodyParts[6]), &result); err != nil {
		return nil, fmt.Errorf("解析数据JSON失败: %w", err)
	}

	return result.Data, nil
}

// QueryStockBasic 查询股票基本信息
func (c *Client) QueryStockBasic(ctx context.Context, code, codeName string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if code != "" {
		if err := validateStockCode(code); err != nil {
			return nil, err
		}
	}

	// 规范化股票代码
	normalizedCode := normalizeStockCode(code)

	msgBody := fmt.Sprintf("query_stock_basic%s%s%s1%s%d%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, normalizedCode, MessageSplit, codeName)

	resp, err := c.sendMessage(ctx, MsgTypeQueryStockBasicRequest, msgBody)
	if err != nil {
		return nil, err
	}

	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 7 {
		return nil, errors.New("无效的响应")
	}

	errorCode := bodyParts[0]
	if errorCode != ErrSuccess {
		errorMsg := ""
		if len(bodyParts) > 1 {
			errorMsg = bodyParts[1]
		}
		return nil, &Error{Code: errorCode, Message: errorMsg}
	}

	// 解析JSON数据
	var result struct {
		Data [][]string `json:"record"`
	}
	if err := json.Unmarshal([]byte(bodyParts[6]), &result); err != nil {
		return nil, fmt.Errorf("解析数据JSON失败: %w", err)
	}

	return result.Data, nil
}

// QueryStockIndustry 查询行业分类
func (c *Client) QueryStockIndustry(ctx context.Context, code, date string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if code != "" {
		if err := validateStockCode(code); err != nil {
			return nil, err
		}
	}

	// 规范化股票代码
	normalizedCode := normalizeStockCode(code)

	msgBody := fmt.Sprintf("query_stock_industry%s%s%s1%s%d%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, normalizedCode, MessageSplit, date)

	resp, err := c.sendMessage(ctx, MsgTypeQueryStockIndustryRequest, msgBody)
	if err != nil {
		return nil, err
	}

	result := &struct {
		ErrorCode string   `json:"-"`
		ErrorMsg  string   `json:"-"`
		Data      [][]string `json:"record"`
	}{}

	if err := parseStandardResponse(resp, result); err != nil {
		return nil, err
	}

	if result.ErrorCode != ErrSuccess {
		return nil, &Error{Code: result.ErrorCode, Message: result.ErrorMsg}
	}

	return result.Data, nil
}

// QueryHS300Stocks 查询沪深300成分股
func (c *Client) QueryHS300Stocks(ctx context.Context, date string) ([][]string, error) {
	return c.queryIndexStocks(ctx, MsgTypeQueryHS300StocksRequest, date)
}

// QuerySZ50Stocks 查询上证50成分股
func (c *Client) QuerySZ50Stocks(ctx context.Context, date string) ([][]string, error) {
	return c.queryIndexStocks(ctx, MsgTypeQuerySZ50StocksRequest, date)
}

// QueryZZ500Stocks 查询中证500成分股
func (c *Client) QueryZZ500Stocks(ctx context.Context, date string) ([][]string, error) {
	return c.queryIndexStocks(ctx, MsgTypeQueryZZ500StocksRequest, date)
}

// queryIndexStocks 指数成分股查询辅助方法
func (c *Client) queryIndexStocks(ctx context.Context, msgType, date string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	msgBody := fmt.Sprintf("index_stocks%s%s%s1%s%d%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, date)

	resp, err := c.sendMessage(ctx, msgType, msgBody)
	if err != nil {
		return nil, err
	}

	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 7 {
		return nil, errors.New("无效的响应")
	}

	errorCode := bodyParts[0]
	if errorCode != ErrSuccess {
		errorMsg := ""
		if len(bodyParts) > 1 {
			errorMsg = bodyParts[1]
		}
		return nil, &Error{Code: errorCode, Message: errorMsg}
	}

	// 解析JSON数据
	var result struct {
		Data [][]string `json:"record"`
	}
	if err := json.Unmarshal([]byte(bodyParts[6]), &result); err != nil {
		return nil, fmt.Errorf("解析数据JSON失败: %w", err)
	}

	return result.Data, nil
}

// QueryDepositRateData 查询存款利率数据
func (c *Client) QueryDepositRateData(ctx context.Context, startDate, endDate string) ([][]string, error) {
	return c.queryEconomicData(ctx, MsgTypeQueryDepositRateDataRequest, startDate, endDate, "")
}

// QueryLoanRateData 查询贷款利率数据
func (c *Client) QueryLoanRateData(ctx context.Context, startDate, endDate string) ([][]string, error) {
	return c.queryEconomicData(ctx, MsgTypeQueryLoanRateDataRequest, startDate, endDate, "")
}

// QueryCPIData 查询CPI数据
func (c *Client) QueryCPIData(ctx context.Context, startDate, endDate string) ([][]string, error) {
	return c.queryEconomicData(ctx, MsgTypeQueryCPIDataRequest, startDate, endDate, "")
}

// QueryPPIData 查询PPI数据
func (c *Client) QueryPPIData(ctx context.Context, startDate, endDate string) ([][]string, error) {
	return c.queryEconomicData(ctx, MsgTypeQueryPPIDataRequest, startDate, endDate, "")
}

// QueryPMIData 查询PMI数据
func (c *Client) QueryPMIData(ctx context.Context, startDate, endDate string) ([][]string, error) {
	return c.queryEconomicData(ctx, MsgTypeQueryPMIDataRequest, startDate, endDate, "")
}

// queryEconomicData 经济数据查询辅助方法
func (c *Client) queryEconomicData(ctx context.Context, msgType, startDate, endDate, extraParam string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	msgBody := fmt.Sprintf("economic_data%s%s%s1%s%d%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, startDate, MessageSplit, endDate)

	if extraParam != "" {
		msgBody += MessageSplit + extraParam
	}

	resp, err := c.sendMessage(ctx, msgType, msgBody)
	if err != nil {
		return nil, err
	}

	result := &struct {
		ErrorCode string   `json:"-"`
		ErrorMsg  string   `json:"-"`
		Data      [][]string `json:"record"`
	}{}

	if err := parseStandardResponse(resp, result); err != nil {
		return nil, err
	}

	if result.ErrorCode != ErrSuccess {
		return nil, &Error{Code: result.ErrorCode, Message: result.ErrorMsg}
	}

	return result.Data, nil
}

// Error 表示 BaoStock 错误
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

// 辅助函数

func formatLength(length int) string {
	return fmt.Sprintf("%010d", length)
}

func validateStockCode(code string) error {
	// 先规范化股票代码（支持6位代码自动添加市场前缀）
	normalizedCode := normalizeStockCode(code)

	// 验证规范化后的代码
	normalizedCode = strings.ToLower(normalizedCode)
	if len(normalizedCode) != StockCodeLength {
		return fmt.Errorf("证券代码必须为%d位，当前: %s", StockCodeLength, code)
	}

	// 检查格式
	if !strings.HasPrefix(normalizedCode, "sh.") && !strings.HasPrefix(normalizedCode, "sz.") {
		return fmt.Errorf("证券代码必须以'sh.'或'sz.'开头，当前: %s", code)
	}

	return nil
}

func parseLoginResponse(resp *Response, result *LoginResponse) error {
	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 4 {
		return errors.New("无效的登录响应")
	}

	result.ErrorCode = bodyParts[0]
	result.ErrorMsg = bodyParts[1]
	result.Method = bodyParts[2]
	result.UserID = bodyParts[3]

	return nil
}

func parseHistoryKDataResponse(resp *Response) (*HistoryKDataResponse, error) {
	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 13 {
		return nil, errors.New("无效的历史K线数据响应")
	}

	result := &HistoryKDataResponse{
		ErrorCode:    bodyParts[0],
		ErrorMsg:     bodyParts[1],
		Method:       bodyParts[2],
		UserID:       bodyParts[3],
		CurPageNum:   bodyParts[4],
		PerPageCount: bodyParts[5],
		Code:         bodyParts[7],
		StartDate:    bodyParts[9],
		EndDate:      bodyParts[10],
		Frequency:    bodyParts[11],
		AdjustFlag:   bodyParts[12],
	}

	// 解析字段
	if len(bodyParts) > 8 {
		fieldsStr := bodyParts[8]
		result.Fields = strings.Split(fieldsStr, ",")
	}

	// 解析JSON数据
	if len(bodyParts) > 6 {
		dataJson := bodyParts[6]
		if dataJson != "" {
			var parsedData struct {
				Record [][]string `json:"record"`
			}
			if err := json.Unmarshal([]byte(dataJson), &parsedData); err != nil {
				return nil, fmt.Errorf("解析数据JSON失败: %w", err)
			}
			result.Data = parsedData.Record
		}
	}

	return result, nil
}

func parseStandardResponse(resp *Response, result interface{}) error {
	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 6 {
		return errors.New("无效的响应")
	}

	// 使用反射或类型断言设置错误字段
	// 为简单起见，假设result有ErrorCode和ErrorMsg字段
	if r, ok := result.(interface{ setErrorCode(string) }); ok {
		r.setErrorCode(bodyParts[0])
	}
	if r, ok := result.(interface{ setErrorMsg(string) }); ok {
		r.setErrorMsg(bodyParts[1])
	}

	// 从bodyParts[6]解析JSON数据
	if len(bodyParts) > 6 && bodyParts[6] != "" {
		if err := json.Unmarshal([]byte(bodyParts[6]), result); err != nil {
			return fmt.Errorf("解析数据JSON失败: %w", err)
		}
	}

	return nil
}
