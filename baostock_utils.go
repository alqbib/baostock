package baostock

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strconv"
	"strings"
)

// calculateCRC32 计算字符串的CRC32校验和
func calculateCRC32(data string) int {
	return int(crc32.ChecksumIEEE([]byte(data)))
}

// decompressData 解压zlib压缩的数据
func decompressData(compressed []byte) ([]byte, error) {
	r, err := zlib.NewReader(strings.NewReader(string(compressed)))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// detectMarketByCode 根据股票代码判断市场
// 上海: 600xxx-604xxx, 688xxx-695xxx (科创板)
// 深圳: 000xxx-004xxx, 300xxx-307xxx (创业板)
func detectMarketByCode(code string) string {
	if len(code) < 3 {
		return "sh"
	}
	prefix := code[0:3]
	switch {
	case prefix >= "600" && prefix <= "604": // 上海主板
		return "sh"
	case prefix >= "688" && prefix <= "695": // 上海科创板
		return "sh"
	case prefix >= "000" && prefix <= "004": // 深圳主板
		return "sz"
	case prefix >= "300" && prefix <= "307": // 深圳创业板
		return "sz"
	default:
		return "sh" // 默认上海
	}
}

// normalizeStockCode 规范化股票代码格式
// 将 "600000.SH" 或 "sh600000" 或 "600000sh" 转换为 "sh.600000"
func normalizeStockCode(code string) string {
	code = strings.ToLower(code)
	code = strings.TrimSpace(code)

	// 已经是正确格式
	if len(code) == StockCodeLength && strings.Contains(code, ".") {
		return code
	}

	// 处理 "600000.sh" 或 "600000SH" 格式
	if len(code) == 9 && !strings.Contains(code, ".") {
		suffix := code[6:8]
		if suffix == "sh" || suffix == "sz" {
			return code[0:6] + "." + code[6:8]
		}
		suffix = code[7:9]
		if suffix == "sh" || suffix == "sz" {
			return code[0:7] + "." + code[7:9]
		}
	}

	// 处理 "sh600000" 格式（缺少点）
	if len(code) == 8 && (strings.HasPrefix(code, "sh") || strings.HasPrefix(code, "sz")) {
		return code[0:2] + "." + code[2:8]
	}

	// 处理 "600000sh" 格式（后缀）
	if len(code) == 8 {
		suffix := code[6:8]
		if suffix == "sh" || suffix == "sz" {
			return suffix + "." + code[0:6]
		}
	}

	// 处理 "600000" 格式 - 根据代码范围自动判断市场
	if len(code) == 6 {
		market := detectMarketByCode(code)
		return market + "." + code
	}

	return code
}

// IsValidStockCode 检查股票代码是否有效
func IsValidStockCode(code string) bool {
	code = strings.ToLower(code)

	if len(code) != StockCodeLength {
		return false
	}

	if !strings.HasPrefix(code, "sh.") && !strings.HasPrefix(code, "sz.") {
		return false
	}

	// 检查代码部分是否全为数字
	codePart := strings.Split(code, ".")[1]
	for _, c := range codePart {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// IsValidDate 检查日期字符串是否为有效格式 (YYYY-MM-DD)
func IsValidDate(date string) bool {
	if len(date) != 10 {
		return false
	}

	if date[4] != '-' || date[7] != '-' {
		return false
	}

	for i := 0; i < 10; i++ {
		if i == 4 || i == 7 {
			continue
		}
		if date[i] < '0' || date[i] > '9' {
			return false
		}
	}

	return true
}

// IsValidYearMonth 检查年月字符串是否为有效格式 (YYYY-MM)
func IsValidYearMonth(ym string) bool {
	if len(ym) != 7 {
		return false
	}

	if ym[4] != '-' {
		return false
	}

	for i := 0; i < 7; i++ {
		if i == 4 {
			continue
		}
		if ym[i] < '0' || ym[i] > '9' {
			return false
		}
	}

	return true
}

// IsValidYear 检查年份字符串是否为有效格式 (YYYY)
func IsValidYear(year string) bool {
	if len(year) != 4 {
		return false
	}

	for _, c := range year {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// ParseMessageHeader 解析消息头字符串
func ParseMessageHeader(header string) (version, msgType string, bodyLength int, err error) {
	parts := strings.Split(header, MessageSplit)
	if len(parts) < 3 {
		err = ErrInvalidHeader
		return
	}

	version = parts[0]
	msgType = parts[1]

	bodyLength, err = strconv.Atoi(parts[2])
	if err != nil {
		err = ErrInvalidHeader
		return
	}

	return
}

// BuildMessageHeader 构建消息头字符串
func BuildMessageHeader(msgType string, bodyLength int) string {
	return fmt.Sprintf("%s%s%s%s%s",
		ClientVersion,
		MessageSplit,
		msgType,
		MessageSplit,
		formatLength(bodyLength))
}

var (
	ErrInvalidHeader = errors.New("无效的消息头")
)

// ToLittleEndianUint32 将uint32转换为小端字节序
func ToLittleEndianUint32(value uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	return buf
}

// FromLittleEndianUint32 将小端字节序转换为uint32
func FromLittleEndianUint32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}
