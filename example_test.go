package baostock

import (
	"context"
	"testing"
	"time"
)

// TestClientLogin 测试客户端登录
func TestClientLogin(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	if !client.loggedIn {
		t.Error("登录后 loggedIn 状态应该为 true")
	}

	if client.userID == "" {
		t.Error("登录后 userID 不应为空")
	}
}

// TestClientQueryHistoryKDataPlus 测试查询历史K线数据
func TestClientQueryHistoryKDataPlus(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryHistoryKDataPlus(context.Background(),
		&HistoryKDataRequest{
			Code:       "sh.600000",
			Fields:     "date,code,open,high,low,close,volume",
			StartDate:  "2023-12-01",
			EndDate:    "2023-12-31",
			Frequency:  FrequencyDaily,
			AdjustFlag: AdjustFlagNoAdjust,
		})
	if err != nil {
		t.Fatalf("查询K线数据失败: %v", err)
	}

	if data.ErrorCode != ErrSuccess {
		t.Fatalf("查询失败: %s", data.ErrorMsg)
	}

	if len(data.Data) == 0 {
		t.Error("应该返回数据")
	}

	if len(data.Fields) != 7 {
		t.Errorf("应该有7个字段, 实际: %d", len(data.Fields))
	}
}

// TestClientQueryTradeDates 测试查询交易日
func TestClientQueryTradeDates(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryTradeDates(context.Background(), "2023-12-01", "2023-12-31")
	if err != nil {
		t.Fatalf("查询交易日失败: %v", err)
	}

	if len(data) == 0 {
		t.Error("应该返回交易日数据")
	}

	// 12月应该有一些交易日
	tradingDays := 0
	for _, row := range data {
		if len(row) > 1 && row[1] == "1" {
			tradingDays++
		}
	}
	if tradingDays == 0 {
		t.Error("12月应该有交易日")
	}
}

// TestClientQueryAllStock 测试查询所有股票
func TestClientQueryAllStock(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryAllStock(context.Background(), "2023-12-29")
	if err != nil {
		t.Fatalf("查询所有股票失败: %v", err)
	}

	if len(data) < 1000 {
		t.Errorf("股票数量太少: %d", len(data))
	}
}

// TestClientQueryHS300Stocks 测试查询沪深300成分股
func TestClientQueryHS300Stocks(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryHS300Stocks(context.Background(), "2023-12-29")
	if err != nil {
		t.Fatalf("查询沪深300失败: %v", err)
	}

	if len(data) != 300 {
		t.Errorf("沪深300应该有300只股票, 实际: %d", len(data))
	}
}

// TestClientQuerySZ50Stocks 测试查询上证50成分股
func TestClientQuerySZ50Stocks(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QuerySZ50Stocks(context.Background(), "2023-12-29")
	if err != nil {
		t.Fatalf("查询上证50失败: %v", err)
	}

	if len(data) != 50 {
		t.Errorf("上证50应该有50只股票, 实际: %d", len(data))
	}
}

// TestClientQueryStockBasic 测试查询股票基本信息
func TestClientQueryStockBasic(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryStockBasic(context.Background(), "sh.600000", "")
	if err != nil {
		t.Fatalf("查询股票基本信息失败: %v", err)
	}

	if len(data) == 0 {
		t.Error("应该返回股票基本信息")
	}
}

// TestClientQueryProfitData 测试查询季频盈利能力
func TestClientQueryProfitData(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryProfitData(context.Background(),
		&QuarterlyDataRequest{
			Code:    "sh.600000",
			Year:    2023,
			Quarter: 4,
		})
	if err != nil {
		t.Fatalf("查询盈利能力数据失败: %v", err)
	}

	if data.ErrorCode != ErrSuccess {
		t.Fatalf("查询失败: %s", data.ErrorMsg)
	}
}

// TestClientQueryDividendData 测试查询股息分红
func TestClientQueryDividendData(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QueryDividendData(context.Background(),
		&DividendDataRequest{
			Code:     "sh.600000",
			Year:     2023,
			YearType: "report",
		})
	if err != nil {
		t.Fatalf("查询股息分红失败: %v", err)
	}

	if data.ErrorCode != ErrSuccess {
		t.Fatalf("查询失败: %s", data.ErrorMsg)
	}
}

// TestClientQuerySTStocks 测试查询ST股票
func TestClientQuerySTStocks(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	data, err := client.QuerySTStocks(context.Background(), "2023-12-29")
	if err != nil {
		t.Fatalf("查询ST股票失败: %v", err)
	}

	// ST股票数量可能为0，所以只检查查询成功
	if data == nil {
		t.Error("应该返回数据")
	}
}

// TestClientWithCustomConfig 测试自定义配置
func TestClientWithCustomConfig(t *testing.T) {
	config := &Config{
		Host:     "www.baostock.com",
		Port:     10030,
		Username: "anonymous",
		Password: "123456",
		Timeout:  30 * time.Second,
	}

	client := NewClientWithConfig(config)

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	if !client.loggedIn {
		t.Error("登录后 loggedIn 状态应该为 true")
	}
}

// Example functions for documentation

func ExampleClient_basicUsage() {
	client := NewClient()

	// 登录
	client.Login(context.Background())
	defer client.Logout(context.Background())

	// 查询历史K线数据
	data, _ := client.QueryHistoryKDataPlus(context.Background(),
		&HistoryKDataRequest{
			Code:       "sh.600000",
			Fields:     "date,code,open,high,low,close,volume",
			StartDate:  "2023-01-01",
			EndDate:    "2023-12-31",
			Frequency:  FrequencyDaily,
			AdjustFlag: AdjustFlagNoAdjust,
		})

	_ = data.ErrorCode
	_ = data.ErrorMsg
	_ = data.Fields
	_ = len(data.Data)
}

func ExampleClient_queryTradeDates() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryTradeDates(context.Background(), "2023-01-01", "2023-01-31")

	for _, row := range data {
		_ = row[0]
		_ = row[1]
	}
}

func ExampleClient_queryAllStock() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryAllStock(context.Background(), "2023-12-29")

	_ = len(data)
	for _, row := range data {
		_, _, _ = row[0], row[1], row[6]
	}
}

func ExampleClient_queryStockBasic() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryStockBasic(context.Background(), "sh.600000", "")

	for _, row := range data {
		_, _, _ = row[0], row[1], row[7]
	}
}

func ExampleClient_queryHS300Stocks() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryHS300Stocks(context.Background(), "2023-12-29")

	_ = len(data)
	for _, row := range data {
		_, _ = row[0], row[1]
	}
}

func ExampleClient_queryDepositRate() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryDepositRateData(context.Background(), "2020-01-01", "2023-12-31")

	for _, row := range data {
		_, _ = row[0], row[1]
	}
}

func ExampleClient_withCustomConfig() {
	config := &Config{
		Host:     "www.baostock.com",
		Port:     10030,
		Username: "anonymous",
		Password: "123456",
		Timeout:  60 * time.Second,
	}

	client := NewClientWithConfig(config)

	client.Login(context.Background())
	defer client.Logout(context.Background())

	// 使用客户端...
	_ = client
}

func Example_errorHandling() {
	client := NewClient()

	err := client.Login(context.Background())
	if err != nil {
		// 检查是否为BaoStock错误
		if bsErr, ok := err.(*Error); ok {
			_ = bsErr.Code
			_ = bsErr.Message
		}
		return
	}
	defer client.Logout(context.Background())

	// 带错误处理的查询
	data, _ := client.QueryHistoryKDataPlus(context.Background(),
		&HistoryKDataRequest{
			Code:       "sh.600000",
			Fields:     "date,code,open,high,low,close,volume",
			StartDate:  "2023-01-01",
			EndDate:    "2023-12-31",
			Frequency:  FrequencyDaily,
			AdjustFlag: AdjustFlagNoAdjust,
		})

	_ = data.ErrorCode
	_ = data.ErrorMsg
	_ = len(data.Data)
}

func ExampleClient_queryProfitData() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryProfitData(context.Background(),
		&QuarterlyDataRequest{
			Code:    "sh.600000",
			Year:    2023,
			Quarter: 4,
		})

	_ = data.ErrorCode
	_ = data.Fields
	for _, row := range data.Data {
		_ = row
	}
}

func ExampleClient_queryDividendData() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QueryDividendData(context.Background(),
		&DividendDataRequest{
			Code:     "sh.600000",
			Year:     2023,
			YearType: "report",
		})

	_ = data.ErrorCode
	for _, row := range data.Data {
		_ = row
	}
}

func ExampleClient_querySTStocks() {
	client := NewClient()

	client.Login(context.Background())
	defer client.Logout(context.Background())

	data, _ := client.QuerySTStocks(context.Background(), "2023-12-29")

	_ = len(data)
	for _, row := range data {
		_, _ = row[0], row[1]
	}
}

// TestErrorString 测试错误类型的字符串输出
func TestErrorString(t *testing.T) {
	err := &Error{Code: ErrSuccess, Message: "success"}
	str := err.Error()
	if str == "" {
		t.Error("错误字符串不应为空")
	}

	err2 := &Error{Code: "99999", Message: "test error"}
	str2 := err2.Error()
	if str2 == "" {
		t.Error("错误字符串不应为空")
	}
}

// TestConstants 测试常量值
func TestConstants(t *testing.T) {
	if DefaultServerHost == "" {
		t.Error("DefaultServerHost 不应为空")
	}
	if DefaultServerPort == 0 {
		t.Error("DefaultServerPort 不应为0")
	}
	if ClientVersion == "" {
		t.Error("ClientVersion 不应为空")
	}
	if DefaultPerPageCount == 0 {
		t.Error("DefaultPerPageCount 不应为0")
	}
}

// TestFrequency 测试频率常量
func TestFrequency(t *testing.T) {
	tests := []struct {
		name     string
		freq     Frequency
		expected string
	}{
		{"5Min", Frequency5Min, "5"},
		{"15Min", Frequency15Min, "15"},
		{"30Min", Frequency30Min, "30"},
		{"60Min", Frequency60Min, "60"},
		{"Daily", FrequencyDaily, "d"},
		{"Weekly", FrequencyWeek, "w"},
		{"Monthly", FrequencyMonth, "m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.freq) != tt.expected {
				t.Errorf("频率 %s 应该是 %s, 实际是 %s", tt.name, tt.expected, tt.freq)
			}
		})
	}
}

// TestAdjustFlag 测试复权标志常量
func TestAdjustFlag(t *testing.T) {
	tests := []struct {
		name     string
		flag     AdjustFlag
		expected string
	}{
		{"Backward", AdjustFlagBackward, "1"},
		{"Forward", AdjustFlagForward, "2"},
		{"NoAdjust", AdjustFlagNoAdjust, "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.flag) != tt.expected {
				t.Errorf("复权标志 %s 应该是 %s, 实际是 %s", tt.name, tt.expected, tt.flag)
			}
		})
	}
}

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Host != DefaultServerHost {
		t.Errorf("默认Host应该是 %s, 实际是 %s", DefaultServerHost, config.Host)
	}
	if config.Port != DefaultServerPort {
		t.Errorf("默认Port应该是 %d, 实际是 %d", DefaultServerPort, config.Port)
	}
	if config.Username != "anonymous" {
		t.Errorf("默认Username应该是 anonymous, 实际是 %s", config.Username)
	}
	if config.Password != "123456" {
		t.Errorf("默认Password应该是 123456, 实际是 %s", config.Password)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("默认Timeout应该是 30秒, 实际是 %v", config.Timeout)
	}
}

// TestNewClient 测试创建客户端
func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("NewClient() 不应返回 nil")
	}

	if client.conn != nil {
		t.Error("新客户端的连接应该为 nil")
	}

	if client.loggedIn {
		t.Error("新客户端的 loggedIn 应该为 false")
	}
}

// TestNewClientWithConfig 测试使用配置创建客户端
func TestNewClientWithConfig(t *testing.T) {
	config := &Config{
		Host:     "testhost",
		Port:     9999,
		Username: "testuser",
		Password: "testpass",
		Timeout:  10 * time.Second,
	}

	client := NewClientWithConfig(config)

	if client == nil {
		t.Fatal("NewClientWithConfig() 不应返回 nil")
	}

	if client.config != config {
		t.Error("客户端配置应该与传入的配置相同")
	}
}

// ExampleNewClient 创建新客户端示例
func ExampleNewClient() {
	client := NewClient()
	_ = client
}

// ExampleNewClientWithConfig 使用自定义配置创建客户端示例
func ExampleNewClientWithConfig() {
	config := &Config{
		Host:     "www.baostock.com",
		Port:     10030,
		Username: "anonymous",
		Password: "123456",
		Timeout:  30 * time.Second,
	}
	client := NewClientWithConfig(config)
	_ = client
}

// TestClientSixDigitStockCode 测试6位股票代码自动规范化
func TestClientSixDigitStockCode(t *testing.T) {
	client := NewClient()

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer client.Logout(context.Background())

	t.Run("上海6位代码", func(t *testing.T) {
		// 测试使用6位代码查询上海股票K线数据
		data, err := client.QueryHistoryKDataPlus(context.Background(),
			&HistoryKDataRequest{
				Code:       "600000", // 6位代码，应自动添加 sh. 前缀
				Fields:     "date,code,open,high,low,close",
				StartDate:  "2023-12-01",
				EndDate:    "2023-12-31",
				Frequency:  FrequencyDaily,
				AdjustFlag: AdjustFlagNoAdjust,
			})
		if err != nil {
			t.Fatalf("查询K线数据失败: %v", err)
		}

		if data.ErrorCode != ErrSuccess {
			t.Fatalf("查询失败: %s", data.ErrorMsg)
		}

		if len(data.Data) == 0 {
			t.Error("应该返回数据")
		}

		// 验证返回的代码是规范化后的格式 (sh.600000)
		if len(data.Data) > 0 && len(data.Data[0]) > 1 {
			returnedCode := data.Data[0][1]
			if returnedCode != "sh.600000" {
				t.Errorf("返回的代码应该是 sh.600000, 实际是: %s", returnedCode)
			}
		}
	})

	t.Run("深圳6位代码", func(t *testing.T) {
		// 测试使用6位代码查询深圳股票K线数据
		data, err := client.QueryHistoryKDataPlus(context.Background(),
			&HistoryKDataRequest{
				Code:       "000001", // 平安银行，应自动添加 sz. 前缀
				Fields:     "date,code,open,high,low,close",
				StartDate:  "2023-12-01",
				EndDate:    "2023-12-31",
				Frequency:  FrequencyDaily,
				AdjustFlag: AdjustFlagNoAdjust,
			})
		if err != nil {
			t.Fatalf("查询K线数据失败: %v", err)
		}

		if data.ErrorCode != ErrSuccess {
			t.Fatalf("查询失败: %s", data.ErrorMsg)
		}

		if len(data.Data) == 0 {
			t.Error("应该返回数据")
		}

		// 验证返回的代码是规范化后的格式 (sz.000001)
		if len(data.Data) > 0 && len(data.Data[0]) > 1 {
			returnedCode := data.Data[0][1]
			if returnedCode != "sz.000001" {
				t.Errorf("返回的代码应该是 sz.000001, 实际是: %s", returnedCode)
			}
		}
	})

	t.Run("创业板6位代码", func(t *testing.T) {
		// 测试使用6位代码查询创业板股票K线数据
		data, err := client.QueryHistoryKDataPlus(context.Background(),
			&HistoryKDataRequest{
				Code:       "300001", // 特锐德，应自动添加 sz. 前缀
				Fields:     "date,code,open,high,low,close",
				StartDate:  "2023-12-01",
				EndDate:    "2023-12-31",
				Frequency:  FrequencyDaily,
				AdjustFlag: AdjustFlagNoAdjust,
			})
		if err != nil {
			t.Fatalf("查询K线数据失败: %v", err)
		}

		if data.ErrorCode != ErrSuccess {
			t.Fatalf("查询失败: %s", data.ErrorMsg)
		}

		if len(data.Data) == 0 {
			t.Error("应该返回数据")
		}

		// 验证返回的代码是规范化后的格式 (sz.300001)
		if len(data.Data) > 0 && len(data.Data[0]) > 1 {
			returnedCode := data.Data[0][1]
			if returnedCode != "sz.300001" {
				t.Errorf("返回的代码应该是 sz.300001, 实际是: %s", returnedCode)
			}
		}
	})
}

// TestNormalizeStockCode 测试股票代码规范化
func TestNormalizeStockCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"6位代码(上海)", "600000", "sh.600000"},
		{"6位代码(深圳主板)", "000001", "sz.000001"},
		{"6位代码(深圳中小板)", "002001", "sz.002001"},
		{"6位代码(深圳创业板)", "300001", "sz.300001"},
		{"6位代码(上海科创板)", "688001", "sh.688001"},
		{"sh.格式", "sh.600000", "sh.600000"},
		{"sz.格式", "sz.000001", "sz.000001"},
		{"sh后缀", "600000sh", "sh.600000"},
		{"无点号", "sh600000", "sh.600000"},
		{"sz后缀", "000001sz", "sz.000001"},
		{"混合大小写", "SH.600000", "sh.600000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeStockCode(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeStockCode(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}
