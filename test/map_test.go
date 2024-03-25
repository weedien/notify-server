package test

import (
	"math/rand"
	"os"
	"runtime"
	"testing"
	"time"
)

const (
	mapSize     = 100_000
	stringLen   = 6
	benchLength = 100_000
)

var m map[string]bool

func TestMain(m *testing.M) {
	createMap(mapSize, stringLen)
	retCode := m.Run()
	os.Exit(retCode)
}

// 10w -> 4MB
// 100w -> 55MB
// 1000w -> 475MB
// 60466176(6**10) -> 3505MB
func BenchmarkMapMemoryUsage(b *testing.B) {
	// 打印内存使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b.Logf("Alloc = %v MB\n", bToMb(m.Alloc))
}

// 创建包含 n 个元素的 map[string]bool
func createMap(n, stringLen int) map[string]bool {
	m = make(map[string]bool, n)
	for i := 0; i < n; i++ {
		key := randomString(stringLen)
		m[key] = true
	}
	return m
}

// 生成指定长度的随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// 执行查询性能测试
func BenchmarkQuery(b *testing.B) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	keys := make([]string, benchLength)
	for i := 0; i < benchLength; i++ {
		keys[i] = randomString(stringLen)
	}

	start := time.Now()
	for _, key := range keys {
		_ = m[key]
	}
	elapsed := time.Since(start)
	b.Logf("Query benchmark (%d ops): %v\n", benchLength, elapsed)
}

// 执行插入性能测试
func BenchmarkInsert(b *testing.B) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	keys := make([]string, benchLength)
	for i := 0; i < benchLength; i++ {
		keys[i] = randomString(stringLen)
	}

	start := time.Now()
	for _, key := range keys {
		m[key] = true
	}
	elapsed := time.Since(start)
	b.Logf("Insert benchmark (%d ops): %v\n", benchLength, elapsed)
}

// 执行随机排序性能测试
func BenchmarkRandomSort(b *testing.B) {
	start := time.Now()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	keys := make([]string, 0, benchLength)
	for code := range m {
		keys = append(keys, code)
	}

	rand.Shuffle(benchLength, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	elapsed := time.Since(start)
	b.Logf("Random sort benchmark (%d ops): %v\n", benchLength, elapsed)
}
