# baostock

[BaoStock](http://baostock.com/) 的 Go 语言客户端库，提供完整的A股历史数据、实时行情和财务数据查询功能。

## 功能特性

- 登录/登出认证
- 历史K线数据查询（支持压缩）
- 交易日查询
- 股票基本信息查询
- 指数成分股（沪深300、上证50、中证500）
- 经济数据（CPI、PPI、PMI、利率等）
- 行业/概念/地域分类
- 季频财务数据（盈利能力、营运能力、成长能力等）
- 股息分红数据
- 复权因子数据
- 完整协议实现

## 安装

```bash
go get github.com/millken/baostock
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    baostock "github.com/millken/baostock"
)

func main() {
    client := baostock.NewClient()

    // 登录
    if err := client.Login(context.Background()); err != nil {
        log.Fatal(err)
    }
    defer client.Logout(context.Background())

    // 查询历史K线数据
    data, err := client.QueryHistoryKDataPlus(context.Background(),
        &baostock.HistoryKDataRequest{
            Code:       "sh.600000",
            Fields:     "date,code,open,high,low,close,volume",
            StartDate:  "2023-01-01",
            EndDate:    "2023-12-31",
            Frequency:  baostock.FrequencyDaily,
            AdjustFlag: baostock.AdjustFlagNoAdjust,
        })
    if err != nil {
        log.Fatal(err)
    }

    for _, row := range data.Data {
        fmt.Printf("%s %s %s\n", row[0], row[1], row[5])
    }
}
```

## 配置选项

```go
config := &baostock.Config{
    Host:     "www.baostock.com",
    Port:     10030,
    Username: "anonymous",
    Password: "123456",
    Timeout:  30 * time.Second,
}
client := baostock.NewClientWithConfig(config)
```

## API 参考

### 客户端方法

| 方法 | 说明 |
|------|------|
| `NewClient()` | 使用默认配置创建客户端 |
| `NewClientWithConfig(config)` | 使用自定义配置创建客户端 |
| `Login(ctx)` | 服务器认证 |
| `Logout(ctx)` | 终止会话 |

### K线数据

```go
data, err := client.QueryHistoryKDataPlus(ctx, &baostock.HistoryKDataRequest{
    Code:       "sh.600000",        // 证券代码 (sh.600000 或 sz.000001)
    Fields:     "date,code,open,high,low,close,volume",
    StartDate:  "2023-01-01",
    EndDate:    "2023-12-31",
    Frequency:  baostock.FrequencyDaily,  // d=日, w=周, m=月, 5/15/30/60=分钟
    AdjustFlag: baostock.AdjustFlagForward, // 1=后复权, 2=前复权, 3=不复权
})
```

### 元数据查询

```go
// 查询交易日
dates, err := client.QueryTradeDates(ctx, "2023-01-01", "2023-12-31")

// 查询所有股票
stocks, err := client.QueryAllStock(ctx, "2023-12-29")

// 查询股票基本信息
info, err := client.QueryStockBasic(ctx, "sh.600000", "")
```

### 指数成分股

```go
// 沪深300
hs300, err := client.QueryHS300Stocks(ctx, "2023-12-29")

// 上证50
sz50, err := client.QuerySZ50Stocks(ctx, "2023-12-29")

// 中证500
zz500, err := client.QueryZZ500Stocks(ctx, "2023-12-29")
```

### 季频财务数据

```go
// 盈利能力
profit, err := client.QueryProfitData(ctx, &baostock.QuarterlyDataRequest{
    Code:    "sh.600000",
    Year:    2023,
    Quarter: 4,
})

// 营运能力
operation, err := client.QueryOperationData(ctx, req)

// 成长能力
growth, err := client.QueryGrowthData(ctx, req)

// 偿债能力
balance, err := client.QueryBalanceData(ctx, req)

// 现金流量
cashFlow, err := client.QueryCashFlowData(ctx, req)

// 杜邦指数
dupont, err := client.QueryDupontData(ctx, req)
```

### 股息分红

```go
data, err := client.QueryDividendData(ctx, &baostock.DividendDataRequest{
    Code:     "sh.600000",
    Year:     2023,
    YearType: "report", // "report"=预案公告年份, "operate"=除权除息年份
})
```

### 复权因子

```go
data, err := client.QueryAdjustFactor(ctx, &baostock.AdjustFactorRequest{
    Code:      "sh.600000",
    StartDate: "2023-01-01",
    EndDate:   "2023-12-31",
})
```

### 经济数据

```go
// 存款利率
rate, err := client.QueryDepositRateData(ctx, "2020-01-01", "2023-12-31")

// 贷款利率
rate, err := client.QueryLoanRateData(ctx, "2020-01-01", "2023-12-31")

// CPI
cpi, err := client.QueryCPIData(ctx, "2020-01-01", "2023-12-31")

// PPI
ppi, err := client.QueryPPIData(ctx, "2020-01-01", "2023-12-31")

// PMI
pmi, err := client.QueryPMIData(ctx, "2020-01-01", "2023-12-31")

// 存款准备金率
ratio, err := client.QueryRequiredReserveRatioData(ctx, "2020-01-01", "2023-12-31", "0")

// 货币供应量（月度）
moneyM, err := client.QueryMoneySupplyDataMonth(ctx, "2020-01", "2023-12")

// 货币供应量（年度）
moneyY, err := client.QueryMoneySupplyDataYear(ctx, "2020", "2023")
```

### 股票分类

```go
// 行业分类
industry, err := client.QueryStockIndustry(ctx, "sh.600000", "2023-12-29")

// 概念分类
concept, err := client.QueryStockConcept(ctx, "sh.600000", "2023-12-29")

// 地域分类
area, err := client.QueryStockArea(ctx, "sh.600000", "2023-12-29")
```

### 特殊股票

```go
// 终止上市股票
terminated, err := client.QueryTerminatedStocks(ctx, "2023-12-29")

// 暂停上市股票
suspended, err := client.QuerySuspendedStocks(ctx, "2023-12-29")

// ST股票
st, err := client.QuerySTStocks(ctx, "2023-12-29")

// *ST股票
starST, err := client.QueryStarSTStocks(ctx, "2023-12-29")

// 中小板
ame, err := client.QueryAMEStocks(ctx, "2023-12-29")

// 创业板
gem, err := client.QueryGEMStocks(ctx, "2023-12-29")

// 沪港通
shhk, err := client.QuerySHHKStocks(ctx, "2023-12-29")

// 深港通
szhk, err := client.QuerySZHKStocks(ctx, "2023-12-29")

// 风险警示板
risk, err := client.QueryStocksInRisk(ctx, "2023-12-29")
```

## 证券代码格式

有效格式：
- `sh.600000` （上海证券交易所）
- `sz.000001` （深圳证券交易所）

## 错误处理

```go
data, err := client.QueryHistoryKDataPlus(ctx, req)
if err != nil {
    // 检查是否为BaoStock错误
    if bsErr, ok := err.(*baostock.Error); ok {
        fmt.Printf("BaoStock错误 [%s]: %s\n", bsErr.Code, bsErr.Message)
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
    return
}

if data.ErrorCode != baostock.ErrSuccess {
    fmt.Printf("服务器错误: %s\n", data.ErrorMsg)
}
```

## 错误代码

| 代码 | 说明 |
|------|------|
| `0` | 成功 |
| `10001001` | 用户未登录 |
| `10001002` | 用户名或密码错误 |
| `10002001` | 网络错误 |
| `10004006` | 参数错误 |
| `10004011` | 无效的证券代码 |

详见 [BAOSTOCK_PROTOCOL.md](BAOSTOCK_PROTOCOL.md) 获取完整错误代码列表。

## 协议文档

本库基于 BaoStock 协议逆向工程实现。详见 [BAOSTOCK_PROTOCOL.md](BAOSTOCK_PROTOCOL.md) 获取完整协议文档，包括：

- 消息格式
- 消息类型代码
- 请求/响应格式
- 数据编码
- 错误代码
- 示例

## 运行示例

```bash
cd examples
go run main.go
```

## 项目结构

```
baostock-go/
├── baostock.go       # 主客户端实现
├── baostock_utils.go # 工具函数
├── baostock_extra.go # 扩展API
├── types.go          # 类型定义
├── example_test.go   # 示例测试
├── examples/
│   └── main.go       # 演示程序
├── go.mod
├── README_GO.md      # 本文档
└── BAOSTOCK_PROTOCOL.md  # 协议文档
```

## 许可证

本项目仅供教育目的使用。请遵守 BaoStock 服务条款。

## 参考资料

- [BaoStock 官方网站](http://baostock.com/)
- [Python BaoStock 库](http://baostock.com/baostock/index.html)

## 贡献

欢迎贡献！请随时提交问题或拉取请求。

## 免责声明

本库与 BaoStock 官方无关联。使用风险自负。
