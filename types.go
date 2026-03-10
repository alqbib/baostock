package baostock

import (
	"fmt"
	"time"
)

// QuarterlyDataRequest 表示季频财务数据请求
type QuarterlyDataRequest struct {
	Code    string // 证券代码
	Year    int    // 年份
	Quarter int    // 季度: 1, 2, 3, 4
}

// QuarterlyDataResponse 表示季频财务数据响应
type QuarterlyDataResponse struct {
	ErrorCode    string   // 错误代码
	ErrorMsg     string   // 错误信息
	Method       string   // 方法名
	UserID       string   // 用户ID
	CurPageNum   string   // 当前页码
	PerPageCount string   // 每页条数
	Data         [][]string // 数据
	Code         string   // 证券代码
	Year         string   // 年份
	Quarter      string   // 季度
	Fields       []string // 字段列表
}

// DividendDataRequest 表示股息分红数据请求
type DividendDataRequest struct {
	Code     string // 证券代码
	Year     int    // 年份
	YearType string // 年份类型: "report"=预案公告年份, "operate"=除权除息年份
}

// AdjustFactorRequest 表示复权因子数据请求
type AdjustFactorRequest struct {
	Code      string // 证券代码
	StartDate string // 开始日期
	EndDate   string // 结束日期
}

// ReportDataRequest 表示公司报告数据请求
type ReportDataRequest struct {
	Code      string // 证券代码
	StartDate string // 开始日期
	EndDate   string // 结束日期
}

// StockClassificationRequest 表示股票分类查询请求
type StockClassificationRequest struct {
	Code string // 证券代码
	Date string // 查询日期
}

// IndexStocksRequest 表示指数成分股查询请求
type IndexStocksRequest struct {
	Date string // 查询日期
}

// EconomicDataRequest 表示经济数据查询请求
type EconomicDataRequest struct {
	StartDate  string // 开始日期
	EndDate    string // 结束日期
	YearType   string // 年份类型（部分指标需要）
}

// DateRange 表示日期范围
type DateRange struct {
	StartDate string // 开始日期
	EndDate   string // 结束日期
}

// NewDateRange 创建新的日期范围
func NewDateRange(start, end string) *DateRange {
	return &DateRange{StartDate: start, EndDate: end}
}

// NewYearDateRange 创建整年的日期范围
func NewYearDateRange(year int) *DateRange {
	return &DateRange{
		StartDate: fmt.Sprintf("%d-01-01", year),
		EndDate:   fmt.Sprintf("%d-12-31", year),
	}
}

// NewMonthDateRange 创建指定月份的日期范围
func NewMonthDateRange(year, month int) *DateRange {
	start := fmt.Sprintf("%d-%02d-01", year, month)
	// 计算月末日期 - 使用下月第0天来获取本月最后一天
	nextMonth := month % 12 // 处理12月的情况，12%12=0，time.Date会正确处理
	nextYear := year
	if nextMonth == 0 {
		nextMonth = 12
	} else {
		nextYear = year + 1
	}
	// 获取下月第一天，然后减去一天得到本月最后一天
	t := time.Date(nextYear, time.Month(nextMonth), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	end := fmt.Sprintf("%d-%02d-%02d", year, month, t.Day())
	return &DateRange{StartDate: start, EndDate: end}
}

// GetCurrentDate 获取当前日期 (YYYY-MM-DD格式)
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurrentYear 获取当前年份
func GetCurrentYear() int {
	return time.Now().Year()
}

// GetCurrentQuarter 获取当前季度 (1-4)
func GetCurrentQuarter() int {
	month := int(time.Now().Month())
	return (month + 2) / 3
}

// ValidateQuarter 检查季度值是否有效
func ValidateQuarter(quarter int) error {
	if quarter < 1 || quarter > 4 {
		return fmt.Errorf("无效的季度: %d，必须为1-4", quarter)
	}
	return nil
}

// ValidateYear 检查年份值是否合理
func ValidateYear(year int) error {
	currentYear := GetCurrentYear()
	if year < 1990 || year > currentYear+1 {
		return fmt.Errorf("无效的年份: %d，必须在1990到%d之间", year, currentYear+1)
	}
	return nil
}
