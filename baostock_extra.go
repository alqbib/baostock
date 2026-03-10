package baostock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// QueryProfitData 查询季频盈利能力数据
func (c *Client) QueryProfitData(ctx context.Context, req *QuarterlyDataRequest) (*QuarterlyDataResponse, error) {
	return c.queryQuarterlyData(ctx, MsgTypeProfitDataRequest, req, "query_profit_data")
}

// QueryOperationData 查询季频营运能力数据
func (c *Client) QueryOperationData(ctx context.Context, req *QuarterlyDataRequest) (*QuarterlyDataResponse, error) {
	return c.queryQuarterlyData(ctx, MsgTypeOperationDataRequest, req, "query_operation_data")
}

// QueryGrowthData 查询季频成长能力数据
func (c *Client) QueryGrowthData(ctx context.Context, req *QuarterlyDataRequest) (*QuarterlyDataResponse, error) {
	return c.queryQuarterlyData(ctx, MsgTypeQueryGrowthDataRequest, req, "query_growth_data")
}

// QueryDupontData 查询季频杜邦指数数据
func (c *Client) QueryDupontData(ctx context.Context, req *QuarterlyDataRequest) (*QuarterlyDataResponse, error) {
	return c.queryQuarterlyData(ctx, MsgTypeQueryDupontDataRequest, req, "query_dupont_data")
}

// QueryBalanceData 查询季频偿债能力数据
func (c *Client) QueryBalanceData(ctx context.Context, req *QuarterlyDataRequest) (*QuarterlyDataResponse, error) {
	return c.queryQuarterlyData(ctx, MsgTypeQueryBalanceDataRequest, req, "query_balance_data")
}

// QueryCashFlowData 查询季频现金流量数据
func (c *Client) QueryCashFlowData(ctx context.Context, req *QuarterlyDataRequest) (*QuarterlyDataResponse, error) {
	return c.queryQuarterlyData(ctx, MsgTypeQueryCashFlowDataRequest, req, "query_cash_flow_data")
}

// queryQuarterlyData 季频数据查询辅助方法
func (c *Client) queryQuarterlyData(ctx context.Context, msgType string, req *QuarterlyDataRequest, methodName string) (*QuarterlyDataResponse, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if err := validateStockCode(req.Code); err != nil {
		return nil, err
	}

	year := req.Year
	if year == 0 {
		year = GetCurrentYear()
	}

	quarter := req.Quarter
	if quarter == 0 {
		quarter = GetCurrentQuarter()
	}

	if err := ValidateQuarter(quarter); err != nil {
		return nil, err
	}

	// 规范化股票代码
	normalizedCode := normalizeStockCode(req.Code)

	msgBody := fmt.Sprintf("%s%s%s%s1%s%d%s%s%s%d%s%d",
		methodName, MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, normalizedCode, MessageSplit, year, MessageSplit, quarter)

	resp, err := c.sendMessage(ctx, msgType, msgBody)
	if err != nil {
		return nil, err
	}

	return parseQuarterlyDataResponse(resp)
}

// QueryDividendData 查询股息分红数据
func (c *Client) QueryDividendData(ctx context.Context, req *DividendDataRequest) (*QuarterlyDataResponse, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if err := validateStockCode(req.Code); err != nil {
		return nil, err
	}

	year := req.Year
	if year == 0 {
		year = GetCurrentYear()
	}

	yearType := req.YearType
	if yearType == "" {
		yearType = "report"
	}

	// 规范化股票代码
	normalizedCode := normalizeStockCode(req.Code)

	msgBody := fmt.Sprintf("query_dividend_data%s%s%s1%s%d%s%s%s%d%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, normalizedCode, MessageSplit, year, MessageSplit, yearType)

	resp, err := c.sendMessage(ctx, MsgTypeQueryDividendDataRequest, msgBody)
	if err != nil {
		return nil, err
	}

	return parseDividendDataResponse(resp)
}

// QueryAdjustFactor 查询复权因子数据
func (c *Client) QueryAdjustFactor(ctx context.Context, req *AdjustFactorRequest) ([][]string, error) {
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
		endDate = GetCurrentDate()
	}

	// 规范化股票代码
	normalizedCode := normalizeStockCode(req.Code)

	msgBody := fmt.Sprintf("query_adjust_factor%s%s%s%s1%d%s%s%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, normalizedCode, MessageSplit, startDate, MessageSplit, endDate)

	resp, err := c.sendMessage(ctx, MsgTypeAdjustFactorRequest, msgBody)
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

// QueryPerformanceExpressReport 查询业绩快报
func (c *Client) QueryPerformanceExpressReport(ctx context.Context, req *ReportDataRequest) ([][]string, error) {
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
		endDate = GetCurrentDate()
	}

	msgBody := fmt.Sprintf("query_performance_express_report%s%s%s%s1%d%s%s%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, req.Code, MessageSplit, startDate, MessageSplit, endDate)

	resp, err := c.sendMessage(ctx, MsgTypeQueryPerformanceExpressReportRequest, msgBody)
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

// QueryForecastReport 查询业绩预告
func (c *Client) QueryForecastReport(ctx context.Context, req *ReportDataRequest) ([][]string, error) {
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
		endDate = GetCurrentDate()
	}

	msgBody := fmt.Sprintf("query_forecast_report%s%s%s%s1%d%s%s%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, req.Code, MessageSplit, startDate, MessageSplit, endDate)

	resp, err := c.sendMessage(ctx, MsgTypeQueryForecastReportRequest, msgBody)
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

// QueryStockConcept 查询概念分类
func (c *Client) QueryStockConcept(ctx context.Context, code, date string) ([][]string, error) {
	return c.queryStockClassification(ctx, MsgTypeQueryStockConceptRequest, code, date, "query_stock_concept")
}

// QueryStockArea 查询地域分类
func (c *Client) QueryStockArea(ctx context.Context, code, date string) ([][]string, error) {
	return c.queryStockClassification(ctx, MsgTypeQueryStockAreaRequest, code, date, "query_stock_area")
}

// queryStockClassification 股票分类查询辅助方法
func (c *Client) queryStockClassification(ctx context.Context, msgType, code, date, methodName string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if code != "" {
		if err := validateStockCode(code); err != nil {
			return nil, err
		}
	}

	msgBody := fmt.Sprintf("%s%s%s%s1%s%d%s%s%s%s",
		methodName, MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, code, MessageSplit, date)

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

// QueryTerminatedStocks 查询终止上市股票
func (c *Client) QueryTerminatedStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQueryTerminatedStocksRequest, "query_terminated_stocks", date)
}

// QuerySuspendedStocks 查询暂停上市股票
func (c *Client) QuerySuspendedStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQuerySuspendedStocksRequest, "query_suspended_stocks", date)
}

// QuerySTStocks 查询ST股票
func (c *Client) QuerySTStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQuerySTStocksRequest, "query_st_stocks", date)
}

// QueryStarSTStocks 查询*ST股票
func (c *Client) QueryStarSTStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQueryStarSTStocksRequest, "query_starst_stocks", date)
}

// QueryAMEStocks 查询中小板股票
func (c *Client) QueryAMEStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQueryAMEStocksRequest, "query_ame_stocks", date)
}

// QueryGEMStocks 查询创业板股票
func (c *Client) QueryGEMStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQueryGEMStocksRequest, "query_gem_stocks", date)
}

// QuerySHHKStocks 查询沪港通股票
func (c *Client) QuerySHHKStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQuerySHHKStocksRequest, "query_shhk_stocks", date)
}

// QuerySZHKStocks 查询深港通股票
func (c *Client) QuerySZHKStocks(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQuerySZHKStocksRequest, "query_szhk_stocks", date)
}

// QueryStocksInRisk 查询风险警示板股票
func (c *Client) QueryStocksInRisk(ctx context.Context, date string) ([][]string, error) {
	return c.querySpecialStocks(ctx, MsgTypeQueryStockInRiskRequest, "query_stocks_in_risk", date)
}

// querySpecialStocks 特殊股票查询辅助方法
func (c *Client) querySpecialStocks(ctx context.Context, msgType, methodName, date string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if date == "" {
		date = GetCurrentDate()
	}

	msgBody := fmt.Sprintf("%s%s%s%s1%s%d%s%s",
		methodName, MessageSplit, c.userID, MessageSplit, MessageSplit,
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

// QueryRequiredReserveRatioData 查询存款准备金率数据
func (c *Client) QueryRequiredReserveRatioData(ctx context.Context, startDate, endDate, yearType string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	if yearType == "" {
		yearType = "0"
	}

	msgBody := fmt.Sprintf("query_required_reserve_ratio_data%s%s%s%s1%d%s%s%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, startDate, MessageSplit, endDate, MessageSplit, yearType)

	resp, err := c.sendMessage(ctx, MsgTypeQueryRequiredReserveRatioDataRequest, msgBody)
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

// QueryMoneySupplyDataMonth 查询月度货币供应量数据
func (c *Client) QueryMoneySupplyDataMonth(ctx context.Context, startDate, endDate string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	msgBody := fmt.Sprintf("query_money_supply_data_month%s%s%s1%s%d%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, startDate, MessageSplit, endDate)

	resp, err := c.sendMessage(ctx, MsgTypeQueryMoneySupplyDataMonthRequest, msgBody)
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

// QueryMoneySupplyDataYear 查询年度货币供应量数据
func (c *Client) QueryMoneySupplyDataYear(ctx context.Context, startDate, endDate string) ([][]string, error) {
	if !c.loggedIn {
		return nil, errors.New("未登录")
	}

	msgBody := fmt.Sprintf("query_money_supply_data_year%s%s%s1%s%d%s%s%s%s",
		MessageSplit, c.userID, MessageSplit, MessageSplit,
		DefaultPerPageCount, MessageSplit, startDate, MessageSplit, endDate)

	resp, err := c.sendMessage(ctx, MsgTypeQueryMoneySupplyDataYearRequest, msgBody)
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

// 解析响应辅助函数

func parseQuarterlyDataResponse(resp *Response) (*QuarterlyDataResponse, error) {
	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 11 {
		return nil, errors.New("无效的季频数据响应")
	}

	result := &QuarterlyDataResponse{
		ErrorCode:    bodyParts[0],
		ErrorMsg:     bodyParts[1],
		Method:       bodyParts[2],
		UserID:       bodyParts[3],
		CurPageNum:   bodyParts[4],
		PerPageCount: bodyParts[5],
		Code:         bodyParts[7],
		Year:         bodyParts[8],
		Quarter:      bodyParts[9],
	}

	// 解析字段
	if len(bodyParts) > 10 {
		fieldsStr := bodyParts[10]
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

func parseDividendDataResponse(resp *Response) (*QuarterlyDataResponse, error) {
	bodyParts := strings.Split(resp.Body, MessageSplit)
	if len(bodyParts) < 11 {
		return nil, errors.New("无效的股息分红数据响应")
	}

	result := &QuarterlyDataResponse{
		ErrorCode:    bodyParts[0],
		ErrorMsg:     bodyParts[1],
		Method:       bodyParts[2],
		UserID:       bodyParts[3],
		CurPageNum:   bodyParts[4],
		PerPageCount: bodyParts[5],
		Code:         bodyParts[7],
		Year:         bodyParts[8],
	}

	// 解析字段
	if len(bodyParts) > 10 {
		fieldsStr := bodyParts[10]
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
