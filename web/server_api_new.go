package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/injoyai/tdx"
	"github.com/injoyai/tdx/protocol"
)

// ==================== 复权/股本变迁接口 ====================

// handleGetGbbq 获取股本变迁/除权除息数据
func handleGetGbbq(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		errorResponse(w, "股票代码不能为空")
		return
	}

	resp, err := client.GetGbbq(code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取股本变迁数据失败: %v", err))
		return
	}

	successResponse(w, resp)
}

// handleGetQfqKline 获取前复权K线数据(基于通达信gbbq,对齐桌面端)
func handleGetQfqKline(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		errorResponse(w, "股票代码不能为空")
		return
	}

	if gbbq == nil {
		errorResponse(w, "复权模块未初始化")
		return
	}

	qfq, err := gbbq.QFQKlineDay(code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取前复权K线失败: %v", err))
		return
	}

	limit := parsePositiveInt(r.URL.Query().Get("limit"))
	list := make([]interface{}, 0, len(qfq))
	for _, k := range qfq {
		list = append(list, k)
		if limit > 0 && len(list) >= limit {
			break
		}
	}

	successResponse(w, map[string]interface{}{
		"count": len(list),
		"list":  list,
		"meta": map[string]interface{}{
			"source": "tdx_gbbq",
			"type":   "day",
			"notes":  []string{"基于通达信股本变迁(gbbq)的前复权数据，对齐通达信桌面端，四舍五入到分"},
		},
	})
}

// handleGetHfqKline 获取后复权K线数据
func handleGetHfqKline(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		errorResponse(w, "股票代码不能为空")
		return
	}

	if gbbq == nil {
		errorResponse(w, "复权模块未初始化")
		return
	}

	hfq, err := gbbq.HFQKlineDay(code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取后复权K线失败: %v", err))
		return
	}

	limit := parsePositiveInt(r.URL.Query().Get("limit"))
	list := make([]interface{}, 0, len(hfq))
	for _, k := range hfq {
		list = append(list, k)
		if limit > 0 && len(list) >= limit {
			break
		}
	}

	successResponse(w, map[string]interface{}{
		"count": len(list),
		"list":  list,
		"meta": map[string]interface{}{
			"source": "tdx_gbbq",
			"type":   "day",
			"notes":  []string{"基于通达信股本变迁(gbbq)的后复权数据，对齐通达信桌面端"},
		},
	})
}

// ==================== 财务/F10接口 ====================

// handleGetFinanceInfo 获取财务信息
func handleGetFinanceInfo(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		errorResponse(w, "股票代码不能为空")
		return
	}

	exchange := determineExchange(code)
	resp, err := client.GetFinanceInfo(exchange, code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取财务信息失败: %v", err))
		return
	}

	successResponse(w, resp)
}

// handleGetF10Category 获取F10公司信息分类目录
func handleGetF10Category(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		errorResponse(w, "股票代码不能为空")
		return
	}

	exchange := determineExchange(code)
	cats, err := client.GetCompanyCategory(exchange, code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取F10分类失败: %v", err))
		return
	}

	successResponse(w, cats)
}

// handleGetF10Content 获取F10某分类的文本内容
func handleGetF10Content(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	filename := strings.TrimSpace(r.URL.Query().Get("filename"))
	startStr := strings.TrimSpace(r.URL.Query().Get("start"))
	lengthStr := strings.TrimSpace(r.URL.Query().Get("length"))

	if code == "" || filename == "" {
		errorResponse(w, "code和filename均为必填参数")
		return
	}

	start := uint32(0)
	if v, err := strconv.ParseUint(startStr, 10, 32); err == nil {
		start = uint32(v)
	}
	length := uint32(10000)
	if v, err := strconv.ParseUint(lengthStr, 10, 32); err == nil && v > 0 {
		length = uint32(v)
	}

	exchange := determineExchange(code)
	content, err := client.GetCompanyContent(exchange, code, filename, start, length)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取F10内容失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"code":     code,
		"filename": filename,
		"start":    start,
		"length":   length,
		"content":  content,
	})
}

// ==================== 板块/行业接口 ====================

// handleGetBlockData 获取板块成分
func handleGetBlockData(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimSpace(r.URL.Query().Get("file"))
	if file == "" {
		file = protocol.BlockFileGN // 默认概念板块
	}

	resp, err := client.GetBlockData(file)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取板块数据失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"file":  file,
		"count": len(resp),
		"list":  resp,
	})
}

// handleGetBlockDataWithIndex 获取板块成分+指数代码ID
func handleGetBlockDataWithIndex(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimSpace(r.URL.Query().Get("file"))
	if file == "" {
		file = protocol.BlockFileGN
	}

	resp, err := client.GetBlockDataWithIndex(file)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取板块数据(含指数代码)失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"file":  file,
		"count": len(resp),
		"list":  resp,
	})
}

// handleGetTdxHy 获取行业归属(通达信/申万)
func handleGetTdxHy(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetTdxHy()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取行业归属失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"count": len(resp),
		"list":  resp,
	})
}

// handleGetTdxZs 获取板块指数代码映射
func handleGetTdxZs(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetTdxZs()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取板块指数代码失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"count": len(resp),
		"list":  resp,
	})
}

// handleGetTdxBk 获取概念板块简称↔全称
func handleGetTdxBk(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetTdxBk()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取概念板块映射失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"count": len(resp),
		"list":  resp,
	})
}

// ==================== 统计/新股接口 ====================

// handleGetTdxStat 获取个股综合统计
func handleGetTdxStat(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetTdxStat()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取个股统计失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"count": len(resp),
		"list":  resp,
	})
}

// handleGetTdxStat2 获取资金流向+板块归属
func handleGetTdxStat2(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetTdxStat2()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取资金流向失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"count": len(resp),
		"list":  resp,
	})
}

// handleGetXgsg 获取新股申购
func handleGetXgsg(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetXgsg()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取新股申购失败: %v", err))
		return
	}

	successResponse(w, map[string]interface{}{
		"count": len(resp),
		"list":  resp,
	})
}

// ==================== 报表/配置接口 ====================

// handleGetZHBFiles 获取zhb.zip内容文件列表
func handleGetZHBFiles(w http.ResponseWriter, r *http.Request) {
	files, err := client.GetZHBFiles()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取zhb.zip失败: %v", err))
		return
	}

	type fileInfo struct {
		Name string `json:"name"`
		Size int    `json:"size"`
	}

	list := make([]fileInfo, 0, len(files))
	for name, data := range files {
		list = append(list, fileInfo{Name: name, Size: len(data)})
	}

	// 如果请求原始文件内容
	download := strings.TrimSpace(r.URL.Query().Get("download"))
	if download != "" {
		data, ok := files[download]
		if !ok {
			errorResponse(w, fmt.Sprintf("文件 %s 不存在于zhb.zip中", download))
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", download))
		encoded := base64.StdEncoding.EncodeToString(data)
		successResponse(w, map[string]interface{}{
			"filename": download,
			"size":     len(data),
			"encoding": "base64",
			"data":     encoded,
		})
		return
	}

	successResponse(w, map[string]interface{}{
		"total_files": len(list),
		"files":       list,
	})
}

// ==================== 扩展行情接口(期货/港股/外盘) ====================

// handleExMarkets 获取扩展行情市场代码表
func handleExMarkets(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	resp, err := exClient.ExMarkets()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取市场代码表失败: %v", err))
		return
	}
	successResponse(w, resp)
}

// handleExCount 获取扩展行情品种数量
func handleExCount(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	resp, err := exClient.ExCount()
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取品种数量失败: %v", err))
		return
	}
	successResponse(w, map[string]interface{}{"count": resp})
}

// handleExInstruments 获取扩展行情品种列表
func handleExInstruments(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	startStr := strings.TrimSpace(r.URL.Query().Get("start"))
	countStr := strings.TrimSpace(r.URL.Query().Get("count"))
	start := uint32(0)
	if v, err := strconv.ParseUint(startStr, 10, 32); err == nil {
		start = uint32(v)
	}
	count := uint16(100)
	if v, err := strconv.ParseUint(countStr, 10, 16); err == nil && v > 0 {
		count = uint16(v)
	}

	resp, err := exClient.ExInstruments(start, count)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取品种列表失败: %v", err))
		return
	}
	successResponse(w, resp)
}

// handleExQuote 获取扩展行情单品种五档
func handleExQuote(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	marketStr := strings.TrimSpace(r.URL.Query().Get("market"))
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if marketStr == "" || code == "" {
		errorResponse(w, "market和code均为必填参数")
		return
	}
	market, err := strconv.ParseUint(marketStr, 10, 8)
	if err != nil {
		errorResponse(w, "market参数无效")
		return
	}

	resp, err := exClient.ExQuote(uint8(market), code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取扩展行情失败: %v", err))
		return
	}
	successResponse(w, resp)
}

// handleExBars 获取扩展行情K线
func handleExBars(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	categoryStr := strings.TrimSpace(r.URL.Query().Get("category"))
	marketStr := strings.TrimSpace(r.URL.Query().Get("market"))
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	startStr := strings.TrimSpace(r.URL.Query().Get("start"))
	countStr := strings.TrimSpace(r.URL.Query().Get("count"))

	if categoryStr == "" || marketStr == "" || code == "" {
		errorResponse(w, "category, market, code均为必填参数")
		return
	}
	category, _ := strconv.ParseUint(categoryStr, 10, 8)
	market, _ := strconv.ParseUint(marketStr, 10, 8)
	start := uint16(0)
	if v, err := strconv.ParseUint(startStr, 10, 16); err == nil {
		start = uint16(v)
	}
	count := uint16(100)
	if v, err := strconv.ParseUint(countStr, 10, 16); err == nil && v > 0 {
		count = uint16(v)
	}

	resp, err := exClient.ExBars(uint8(category), uint8(market), code, start, count)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取扩展K线失败: %v", err))
		return
	}
	successResponse(w, resp)
}

// handleExMinute 获取扩展行情当日分时
func handleExMinute(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	marketStr := strings.TrimSpace(r.URL.Query().Get("market"))
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if marketStr == "" || code == "" {
		errorResponse(w, "market和code均为必填参数")
		return
	}
	market, _ := strconv.ParseUint(marketStr, 10, 8)

	resp, err := exClient.ExMinute(uint8(market), code)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取扩展分时失败: %v", err))
		return
	}
	successResponse(w, resp)
}

// handleExTrade 获取扩展行情分笔成交
func handleExTrade(w http.ResponseWriter, r *http.Request) {
	if exClient == nil {
		errorResponse(w, "扩展行情客户端未连接")
		return
	}
	marketStr := strings.TrimSpace(r.URL.Query().Get("market"))
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	startStr := strings.TrimSpace(r.URL.Query().Get("start"))
	countStr := strings.TrimSpace(r.URL.Query().Get("count"))

	if marketStr == "" || code == "" {
		errorResponse(w, "market和code均为必填参数")
		return
	}
	market, _ := strconv.ParseUint(marketStr, 10, 8)
	start := uint16(0)
	if v, err := strconv.ParseUint(startStr, 10, 16); err == nil {
		start = uint16(v)
	}
	count := uint16(100)
	if v, err := strconv.ParseUint(countStr, 10, 16); err == nil && v > 0 {
		count = uint16(v)
	}

	resp, err := exClient.ExTrade(uint8(market), code, start, count)
	if err != nil {
		errorResponse(w, fmt.Sprintf("获取扩展分笔成交失败: %v", err))
		return
	}
	successResponse(w, resp)
}

// ==================== 辅助函数 ====================

// determineExchange 根据股票代码判断交易所
func determineExchange(code string) protocol.Exchange {
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "sh") || strings.HasPrefix(code, "SH") {
		return protocol.ExchangeSH
	}
	if strings.HasPrefix(code, "sz") || strings.HasPrefix(code, "SZ") {
		return protocol.ExchangeSZ
	}
	if strings.HasPrefix(code, "bj") || strings.HasPrefix(code, "BJ") {
		return protocol.ExchangeBJ
	}
	// 根据代码数字判断
	switch {
	case strings.HasPrefix(code, "6"), strings.HasPrefix(code, "9"):
		return protocol.ExchangeSH
	case strings.HasPrefix(code, "0"), strings.HasPrefix(code, "1"), strings.HasPrefix(code, "2"), strings.HasPrefix(code, "3"):
		return protocol.ExchangeSZ
	case strings.HasPrefix(code, "4"), strings.HasPrefix(code, "8"):
		return protocol.ExchangeBJ
	default:
		return protocol.ExchangeSH
	}
}

// ensure DefaultCodes helper still works with new tdx API
func _() {
	_ = tdx.DefaultCodes
	_ = tdx.DialExHqDefault
}
