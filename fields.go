package baostock

// K线数据字段常量
// 参考: https://www.baostock.com/mainContent?file=stockKData.md

const (
	// 通用字段
	FieldDate       = "date"       // 交易所行情日期，格式：YYYY-MM-DD
	FieldCode       = "code"       // 证券代码，格式：sh.600000 或 sz.000001
	FieldOpen       = "open"       // 开盘价，精度：小数点后4位，单位：人民币元
	FieldHigh       = "high"       // 最高价，精度：小数点后4位，单位：人民币元
	FieldLow        = "low"        // 最低价，精度：小数点后4位，单位：人民币元
	FieldClose      = "close"      // 收盘价，精度：小数点后4位，单位：人民币元
	FieldVolume     = "volume"     // 成交数量，单位：股
	FieldAmount     = "amount"     // 成交金额，精度：小数点后4位，单位：人民币元
	FieldAdjustFlag = "adjustflag" // 复权状态：1=后复权，2=前复权，3=不复权

	// 日线专用字段
	FieldPreClose   = "preclose"    // 昨日收盘价，精度：小数点后4位，单位：人民币元
	FieldTurn       = "turn"        // 换手率，精度：小数点后6位，单位：%
	FieldTradeStatus = "tradestatus" // 交易状态：1=正常交易，0=停牌
	FieldPctChg     = "pctChg"      // 涨跌幅（百分比），精度：小数点后6位
	FieldPETTM      = "peTTM"       // 滚动市盈率，精度：小数点后6位
	FieldPSTTM      = "psTTM"       // 滚动市销率，精度：小数点后6位
	FieldPcfNcfTTM  = "pcfNcfTTM"   // 滚动市现率，精度：小数点后6位
	FieldPbMRQ      = "pbMRQ"       // 市净率，精度：小数点后6位
	FieldIsST       = "isST"        // 是否ST股：1=是，0=否

	// 分钟线专用字段
	FieldTime = "time" // 交易所行情时间，格式：YYYYMMDDHHMMSSsss
)

// DailyKLineFields 日线字段集合（包含停牌证券）
var DailyKLineFields = []string{
	FieldDate,
	FieldCode,
	FieldOpen,
	FieldHigh,
	FieldLow,
	FieldClose,
	FieldPreClose,
	FieldVolume,
	FieldAmount,
	FieldAdjustFlag,
	FieldTurn,
	FieldTradeStatus,
	FieldPctChg,
	FieldPETTM,
	FieldPSTTM,
	FieldPcfNcfTTM,
	FieldPbMRQ,
	FieldIsST,
}

// DailyKLineCommonFields 日线常用字段集合
var DailyKLineCommonFields = []string{
	FieldDate,
	FieldCode,
	FieldOpen,
	FieldHigh,
	FieldLow,
	FieldClose,
	FieldVolume,
	FieldAmount,
}

// WeeklyMonthlyKLineFields 周月线字段集合
var WeeklyMonthlyKLineFields = []string{
	FieldDate,
	FieldCode,
	FieldOpen,
	FieldHigh,
	FieldLow,
	FieldClose,
	FieldVolume,
	FieldAmount,
	FieldAdjustFlag,
	FieldTurn,
	FieldPctChg,
}

// MinuteKLineFields 分钟线字段集合（5、15、30、60分钟）
var MinuteKLineFields = []string{
	FieldDate,
	FieldTime,
	FieldCode,
	FieldOpen,
	FieldHigh,
	FieldLow,
	FieldClose,
	FieldVolume,
	FieldAmount,
	FieldAdjustFlag,
}

