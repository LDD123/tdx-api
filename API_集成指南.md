# 📡 API功能完整集成指南

## 🎯 概述

已集成通达信协议库(tdx)最新版本的全部功能，包括基础行情、扩展数据、复权计算、财务信息、板块数据、统计指标、扩展行情(期货/港股)等。所有接口已默认集成，开箱即用。

### ✅ 已实现的基础接口（6个）
1. **GET /api/quote** - 五档行情
2. **GET /api/kline** - K线数据
3. **GET /api/minute** - 分时数据
4. **GET /api/trade** - 分时成交
5. **GET /api/search** - 搜索股票
6. **GET /api/stock-info** - 综合信息

### ✅ 扩展接口（7个）
7. **GET /api/codes** - 股票代码列表
8. **POST /api/batch-quote** - 批量获取行情
9. **GET /api/kline-history** - 历史K线范围查询
10. **GET /api/index** - 指数数据
11. **GET /api/market-stats** - 市场统计
12. **GET /api/server-status** - 服务状态
13. **GET /api/health** - 健康检查

### ✅ 数据入库任务接口（5个）
14. **POST /api/tasks/pull-kline** - 批量K线入库任务
15. **POST /api/tasks/pull-trade** - 分时成交入库任务
16. **GET /api/tasks** - 查询任务列表
17. **GET /api/tasks/{id}** - 查询任务详情
18. **POST /api/tasks/{id}/cancel** - 取消任务

### ✅ 数据服务接口（13个）
19. **GET /api/etf** - ETF基金列表
20. **GET /api/trade-history** - 历史分时成交分页
21. **GET /api/minute-trade-all** - 全天分时成交汇总
22. **GET /api/workday** - 交易日信息查询
23. **GET /api/market-count** - 各交易所证券数量
24. **GET /api/stock-codes** - 全部股票代码
25. **GET /api/etf-codes** - 全部ETF代码
26. **GET /api/kline-all** - 股票历史K线全集
27. **GET /api/index/all** - 指数历史K线全集
28. **GET /api/trade-history/full** - 上市以来分时成交
29. **GET /api/workday/range** - 交易日范围列表
30. **GET /api/income** - 收益区间分析
31. **GET /api/call-auction** - 集合竞价数据

### ✅ 全量历史K线接口（2个）
32. **GET /api/kline-all/tdx** - 通达信原始历史K线
33. **GET /api/kline-all/ths** - 同花顺前复权历史K线

### ✅ 新增：复权/股本变迁接口（3个）🆕
34. **GET /api/gbbq** - 获取股本变迁/除权除息数据
35. **GET /api/qfq-kline** - 获取前复权K线（基于通达信gbbq，对齐桌面端）
36. **GET /api/hfq-kline** - 获取后复权K线

### ✅ 新增：财务/F10接口（3个）🆕
37. **GET /api/finance** - 获取财务信息（流通股本/总股本/净利润等）
38. **GET /api/f10/category** - 获取F10公司信息分类目录
39. **GET /api/f10/content** - 获取F10某分类的文本内容

### ✅ 新增：板块/行业接口（5个）🆕
40. **GET /api/block** - 获取板块成分（概念/地域/风格）
41. **GET /api/block-with-index** - 获取板块成分+指数代码ID
42. **GET /api/tdxhy** - 获取行业归属（通达信/申万）
43. **GET /api/tdxzs** - 获取板块指数代码映射
44. **GET /api/tdxbk** - 获取概念板块简称↔全称

### ✅ 新增：统计/新股接口（3个）🆕
45. **GET /api/tdxstat** - 获取个股综合统计（市盈率/股息率/涨跌幅等）
46. **GET /api/tdxstat2** - 获取资金流向+板块归属
47. **GET /api/xgsg** - 获取新股申购列表

### ✅ 新增：报表/配置接口（1个）🆕
48. **GET /api/zhb-files** - 获取zhb.zip内容文件列表

### ✅ 新增：扩展行情接口（7个）🆕
49. **GET /api/ex/markets** - 扩展行情市场代码表
50. **GET /api/ex/count** - 扩展行情品种数量
51. **GET /api/ex/quote** - 扩展行情单品种五档
52. **GET /api/ex/bars** - 扩展行情K线
53. **GET /api/ex/minute** - 扩展行情当日分时
54. **GET /api/ex/trade** - 扩展行情分笔成交

---

## 🆕 新增功能详解

### 1. 复权/股本变迁（GBBQ）

基于通达信股本变迁(gbbq)数据，实现与通达信桌面端对齐的前/后复权计算。

**核心特性**：
- 使用仿射变换模型 `price_adj = QFQMul × price_raw + QFQAdd`，四舍五入到分
- 与通达信桌面端逐日对齐（含配股、大比例送转、股改停牌复合事件）
- 自动缓存gbbq数据到本地SQLite，定时更新
- 替代原有的同花顺爬虫方式，更稳定可靠

```bash
# 获取前复权日K线（推荐，对齐通达信）
curl "http://localhost:8080/api/qfq-kline?code=000001"

# 获取后复权日K线
curl "http://localhost:8080/api/hfq-kline?code=000001"

# 获取除权除息数据
curl "http://localhost:8080/api/gbbq?code=000001"
```

### 2. 财务信息

获取标的的财务/基本面数据，包括流通股本、总股本、上市日期、股东户数、净利润等30+字段。

```bash
# 获取财务信息
curl "http://localhost:8080/api/finance?code=600519"

# 获取F10分类目录
curl "http://localhost:8080/api/f10/category?code=600519"

# 获取F10内容
curl "http://localhost:8080/api/f10/content?code=600519&filename=company.txt&start=0&length=5000"
```

### 3. 板块/行业数据

支持概念板块、地域风格、指数板块的成分查询，以及通达信/申万行业归属。

```bash
# 获取概念板块成分（默认）
curl "http://localhost:8080/api/block"

# 获取风格板块（含地域）
curl "http://localhost:8080/api/block?file=block_fg.dat"

# 获取板块+指数代码ID（自动关联tdxzs.cfg）
curl "http://localhost:8080/api/block-with-index?file=block_gn.dat"

# 获取行业归属
curl "http://localhost:8080/api/tdxhy"

# 获取板块指数代码映射
curl "http://localhost:8080/api/tdxzs"

# 获取概念板块简称↔全称
curl "http://localhost:8080/api/tdxbk"
```

**板块文件说明**：
| 文件名 | 说明 |
|--------|------|
| `block_gn.dat` | 概念板块（默认） |
| `block_fg.dat` | 风格板块（含地域） |
| `block_zs.dat` | 指数板块 |
| `block_hy.dat` | 行业板块 |

### 4. 个股统计/新股申购

来自通达信zhb.zip的全市场逐股数据，经10只大市值股对照实盘核验。

```bash
# 获取个股综合统计（市盈率TTM/静态市盈率/股息率/涨跌幅/连涨连跌/区间涨跌幅）
curl "http://localhost:8080/api/tdxstat"

# 获取资金流向+板块归属
curl "http://localhost:8080/api/tdxstat2"

# 获取新股申购
curl "http://localhost:8080/api/xgsg"
```

### 5. 扩展行情（期货/港股/外盘）

通过独立端口(7727)连接扩展行情服务器，获取期货、港股、外盘等数据。

```bash
# 获取市场代码表
curl "http://localhost:8080/api/ex/markets"

# 获取期货五档行情
curl "http://localhost:8080/api/ex/quote?market=29&code=IF2401"

# 获取港股五档行情
curl "http://localhost:8080/api/ex/quote?market=31&code=00700"

# 获取K线数据
curl "http://localhost:8080/api/ex/bars?category=4&market=29&code=IF2401&start=0&count=100"

# 获取分时数据
curl "http://localhost:8080/api/ex/minute?market=31&code=00700"

# 获取分笔成交
curl "http://localhost:8080/api/ex/trade?market=31&code=00700&start=0&count=50"
```

---

## 🚀 如何集成

> 当前仓库已经完成以下步骤，接口可直接使用；若需要迁移到其他工程或自定义修改，可参考下述说明。

### 路由注册（server.go main函数）

```go
func main() {
    // 静态文件服务
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // === 现有API路由 ===
    http.HandleFunc("/api/quote", handleGetQuote)
    http.HandleFunc("/api/kline", handleGetKline)
    http.HandleFunc("/api/minute", handleGetMinute)
    http.HandleFunc("/api/trade", handleGetTrade)
    http.HandleFunc("/api/search", handleSearchCode)
    http.HandleFunc("/api/stock-info", handleGetStockInfo)

    // === 扩展API路由 ===
    http.HandleFunc("/api/codes", handleGetCodes)
    http.HandleFunc("/api/batch-quote", handleBatchQuote)
    http.HandleFunc("/api/kline-history", handleGetKlineHistory)
    http.HandleFunc("/api/index", handleGetIndex)
    http.HandleFunc("/api/index/all", handleGetIndexAll)
    http.HandleFunc("/api/market-stats", handleGetMarketStats)
    http.HandleFunc("/api/market-count", handleGetMarketCount)
    http.HandleFunc("/api/stock-codes", handleGetStockCodes)
    http.HandleFunc("/api/etf-codes", handleGetETFCodes)
    http.HandleFunc("/api/server-status", handleGetServerStatus)
    http.HandleFunc("/api/health", handleHealthCheck)
    http.HandleFunc("/api/etf", handleGetETFList)
    http.HandleFunc("/api/trade-history", handleGetTradeHistory)
    http.HandleFunc("/api/trade-history/full", handleGetTradeHistoryFull)
    http.HandleFunc("/api/minute-trade-all", handleGetMinuteTradeAll)
    http.HandleFunc("/api/kline-all", handleGetKlineAllTDX)
    http.HandleFunc("/api/kline-all/tdx", handleGetKlineAllTDX)
    http.HandleFunc("/api/kline-all/ths", handleGetKlineAllTHS)
    http.HandleFunc("/api/workday", handleGetWorkday)
    http.HandleFunc("/api/workday/range", handleGetWorkdayRange)
    http.HandleFunc("/api/income", handleGetIncome)
    http.HandleFunc("/api/tasks/pull-kline", handleCreatePullKlineTask)
    http.HandleFunc("/api/tasks/pull-trade", handleCreatePullTradeTask)
    http.HandleFunc("/api/tasks", handleListTasks)
    http.HandleFunc("/api/tasks/", handleTaskOperations)
    http.HandleFunc("/api/call-auction", handleGetCallAuction)

    // === 新增：复权/股本变迁接口 ===
    http.HandleFunc("/api/gbbq", handleGetGbbq)
    http.HandleFunc("/api/qfq-kline", handleGetQfqKline)
    http.HandleFunc("/api/hfq-kline", handleGetHfqKline)

    // === 新增：财务/F10接口 ===
    http.HandleFunc("/api/finance", handleGetFinanceInfo)
    http.HandleFunc("/api/f10/category", handleGetF10Category)
    http.HandleFunc("/api/f10/content", handleGetF10Content)

    // === 新增：板块/行业接口 ===
    http.HandleFunc("/api/block", handleGetBlockData)
    http.HandleFunc("/api/block-with-index", handleGetBlockDataWithIndex)
    http.HandleFunc("/api/tdxhy", handleGetTdxHy)
    http.HandleFunc("/api/tdxzs", handleGetTdxZs)
    http.HandleFunc("/api/tdxbk", handleGetTdxBk)

    // === 新增：统计/新股接口 ===
    http.HandleFunc("/api/tdxstat", handleGetTdxStat)
    http.HandleFunc("/api/tdxstat2", handleGetTdxStat2)
    http.HandleFunc("/api/xgsg", handleGetXgsg)

    // === 新增：报表/配置接口 ===
    http.HandleFunc("/api/zhb-files", handleGetZHBFiles)

    // === 新增：扩展行情接口(期货/港股/外盘) ===
    http.HandleFunc("/api/ex/markets", handleExMarkets)
    http.HandleFunc("/api/ex/count", handleExCount)
    http.HandleFunc("/api/ex/quote", handleExQuote)
    http.HandleFunc("/api/ex/bars", handleExBars)
    http.HandleFunc("/api/ex/minute", handleExMinute)
    http.HandleFunc("/api/ex/trade", handleExTrade)

    port := ":8080"
    log.Printf("服务启动成功，访问 http://localhost%s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
```

### 初始化代码（server.go init函数）

```go
func init() {
    // 连接通达信服务器
    client, err = tdx.DialDefault(tdx.WithDebug(false))
    if err != nil {
        log.Fatalf("连接服务器失败: %v", err)
    }

    // 初始化代码缓存
    if codes, err := tdx.NewCodesSqlite(client); err != nil {
        log.Printf("初始化代码库失败: %v", err)
    } else {
        tdx.DefaultCodes = codes
        tdx.DefaultCodes.Update()
    }

    // 初始化数据管理器（新版使用Option模式）
    manager, err = tdx.NewManage(tdx.WithClients(4))
    if err != nil {
        log.Fatalf("初始化数据管理器失败: %v", err)
    }
    manager.Codes.Update()
    manager.Workday.Update()
    manager.Cron.Start()

    // 初始化复权模块（基于通达信gbbq，对齐桌面端）
    if g, err := tdx.NewGbbq(tdx.WithGbbqClient(client)); err != nil {
        log.Printf("初始化复权模块失败: %v", err)
    } else {
        gbbq = g
    }

    // 初始化扩展行情客户端（期货/港股/外盘，端口7727，可选）
    if ec, err := tdx.DialExHqDefault(tdx.WithDebug(false)); err != nil {
        log.Printf("连接扩展行情服务器失败(可选): %v", err)
    } else {
        exClient = ec
    }
}
```

---

## ⚠️ 版本更新说明（v1.x → v2.x）

### 破坏性变更

| 旧接口/方法 | 新接口/方法 | 说明 |
|------------|------------|------|
| `tdx.NewManage(&tdx.ManageConfig{Number: N})` | `tdx.NewManage(tdx.WithClients(N))` | ManageConfig 改为 Option 模式 |
| `extend.KlineTableMap` | `extend.Day` / `extend.Minute` | 表映射改为常量 |
| `extend.PullKlineConfig.Tables` | `extend.PullKlineConfig.Types` | 字段重命名 |
| `extend.PullKlineConfig.Limit` | `extend.PullKlineConfig.Goroutines` | 字段重命名 |
| `extend.PullTrade.StartYear/EndYear` | 已移除 | 固定从2000年拉取 |
| `extend.Kline{Code, Date, ...}` | `extend.Kline{Unix, *protocol.Kline, ...}` | 结构体字段变更 |

### 新增依赖

| 依赖 | 版本 | 说明 |
|------|------|------|
| `github.com/injoyai/bar` | v0.0.9 | 进度条 |
| `github.com/injoyai/base` | v1.2.20 | 基础库（从v1.2.17升级） |
| `github.com/grafov/m3u8` | v0.12.1 | M3U8解析 |

---

## 🧪 测试新接口

### 测试1: 前复权K线

```bash
curl "http://localhost:8080/api/qfq-kline?code=000001&limit=10"
```

### 测试2: 财务信息

```bash
curl "http://localhost:8080/api/finance?code=600519"
```

### 测试3: 概念板块

```bash
curl "http://localhost:8080/api/block-with-index?file=block_gn.dat"
```

### 测试4: 个股统计

```bash
curl "http://localhost:8080/api/tdxstat"
```

### 测试5: 扩展行情

```bash
curl "http://localhost:8080/api/ex/markets"
curl "http://localhost:8080/api/ex/quote?market=31&code=00700"
```

---

## 📚 完整API列表

### 基础数据接口

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/quote | GET | 五档行情 |
| /api/kline | GET | K线数据（含日/周/月前复权） |
| /api/minute | GET | 分时数据 |
| /api/trade | GET | 分时成交 |
| /api/search | GET | 搜索股票 |
| /api/stock-info | GET | 综合信息汇总 |

### 扩展功能接口

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/codes | GET | 股票列表 |
| /api/batch-quote | POST | 批量行情 |
| /api/kline-history | GET | 历史K线 |
| /api/index | GET | 指数数据 |
| /api/market-stats | GET | 市场统计 |
| /api/market-count | GET | 市场数量 |
| /api/stock-codes | GET | 股票代码 |
| /api/etf-codes | GET | ETF代码 |
| /api/server-status | GET | 服务状态 |
| /api/health | GET | 健康检查 |
| /api/etf | GET | ETF列表 |
| /api/trade-history | GET | 历史成交 |
| /api/trade-history/full | GET | 完整历史成交 |
| /api/minute-trade-all | GET | 全天分时成交 |
| /api/kline-all | GET | K线全集 |
| /api/kline-all/tdx | GET | TDX源K线 |
| /api/kline-all/ths | GET | 同花顺源K线 |
| /api/workday | GET | 交易日查询 |
| /api/workday/range | GET | 交易日范围 |
| /api/income | GET | 收益分析 |
| /api/call-auction | GET | 集合竞价 |

### 复权/股本变迁接口 🆕

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/gbbq | GET | 股本变迁/除权除息 |
| /api/qfq-kline | GET | 前复权K线（对齐通达信桌面端） |
| /api/hfq-kline | GET | 后复权K线 |

### 财务/F10接口 🆕

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/finance | GET | 财务信息 |
| /api/f10/category | GET | F10分类目录 |
| /api/f10/content | GET | F10内容 |

### 板块/行业接口 🆕

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/block | GET | 板块成分 |
| /api/block-with-index | GET | 板块成分+指数代码 |
| /api/tdxhy | GET | 行业归属 |
| /api/tdxzs | GET | 板块指数代码映射 |
| /api/tdxbk | GET | 概念板块简称↔全称 |

### 统计/新股接口 🆕

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/tdxstat | GET | 个股综合统计 |
| /api/tdxstat2 | GET | 资金流向+板块归属 |
| /api/xgsg | GET | 新股申购 |

### 报表/配置接口 🆕

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/zhb-files | GET | zhb.zip文件列表 |

### 扩展行情接口 🆕

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/ex/markets | GET | 扩展行情市场代码表 |
| /api/ex/count | GET | 扩展行情品种数量 |
| /api/ex/quote | GET | 扩展行情五档 |
| /api/ex/bars | GET | 扩展行情K线 |
| /api/ex/minute | GET | 扩展行情分时 |
| /api/ex/trade | GET | 扩展行情分笔成交 |

### 数据入库接口

| 接口 | 方法 | 说明 |
|-----|------|------|
| /api/tasks/pull-kline | POST | K线入库任务 |
| /api/tasks/pull-trade | POST | 成交入库任务 |
| /api/tasks | GET | 任务列表 |
| /api/tasks/{id} | GET | 任务详情 |
| /api/tasks/{id}/cancel | POST | 取消任务 |

---

## ✅ 总结

### 已完成
✅ 54个完整API接口
✅ 基于通达信gbbq的前/后复权（对齐桌面端）
✅ 财务信息/F10公司资料
✅ 板块成分/行业归属
✅ 个股统计/资金流向/新股申购
✅ 扩展行情（期货/港股/外盘）
✅ 详细的接口文档
✅ 使用示例
✅ 集成指南

### 使用流程
1. 阅读 `API_接口文档.md` 了解所有接口
2. 按照本文档集成扩展接口
3. 重新构建Docker镜像
4. 测试接口功能
5. 开始使用API开发应用
