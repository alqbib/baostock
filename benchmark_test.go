package baostock

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 生成 10000 条记录的 JSON 数据
func generateLargeJsonData(count int) string {
	records := make([][]string, count)
	for i := 0; i < count; i++ {
		records[i] = []string{"2023-12-01", "sh.600000", "6.85", "6.90", "6.80", "6.85", "100000000", "685000000.00", "1"}
	}

	data := struct {
		Record [][]string `json:"record"`
	}{
		Record: records,
	}

	jsonBytes, _ := json.Marshal(data)
	return string(jsonBytes)
}

var mockJsonData = generateLargeJsonData(10000)

// BenchmarkStreamJsonRecords 流式解析性能测试
func BenchmarkStreamJsonRecords(b *testing.B) {
	b.ResetTimer()
	recordCount := 0

	err := streamJsonRecords(mockJsonData, func(record []string) error {
		recordCount++
		return nil
	})

	if err != nil {
		b.Fatalf("流式解析失败: %v", err)
	}

	fmt.Printf("\n流式解析处理了 %d 条记录\n", recordCount)
}

// BenchmarkParseJsonRecords 一次性解析性能测试
func BenchmarkParseJsonRecords(b *testing.B) {
	b.ResetTimer()

	records, err := parseJsonRecords(mockJsonData)
	if err != nil {
		b.Fatalf("一次性解析失败: %v", err)
	}

	fmt.Printf("\n一次性解析处理了 %d 条记录\n", len(records))
}

// BenchmarkStreamJsonRecordsMemory 流式解析内存测试
func BenchmarkStreamJsonRecordsMemory(b *testing.B) {
	b.ReportAllocs()
	recordCount := 0

	err := streamJsonRecords(mockJsonData, func(record []string) error {
		recordCount++
		return nil
	})

	if err != nil {
		b.Fatalf("流式解析失败: %v", err)
	}
}

// BenchmarkParseJsonRecordsMemory 一次性解析内存测试
func BenchmarkParseJsonRecordsMemory(b *testing.B) {
	b.ReportAllocs()

	_, err := parseJsonRecords(mockJsonData)
	if err != nil {
		b.Fatalf("一次性解析失败: %v", err)
	}
}

// BenchmarkIterationCompare 迭代方式对比
func BenchmarkIterationCompare(b *testing.B) {
	// 模拟 10000 条记录
	records := make([][]string, 10000)
	for i := 0; i < 10000; i++ {
		records[i] = []string{"2023-12-01", "sh.600000", "6.85"}
	}

	b.Run("CallbackStyle", func(b *testing.B) {
		b.ReportAllocs()
		count := 0
		for _, record := range records {
			count++
			_ = record
		}
	})

	b.Run("DirectLoop", func(b *testing.B) {
		b.ReportAllocs()
		count := 0
		for _, record := range records {
			count++
			_ = record
		}
	})
}
