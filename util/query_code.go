package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/willf/bloom"
)

var QueryCodeLen = 4

// var queryCodeSet = make(map[string]bool)

// var codeCollectPeriod = 24 * time.Hour

const n = uint(1000)

var filter = bloom.New(20*n, 5)

// 生成4次查询码，如果没有生成到不重复的查询码，就返回一个过期的查询码
func GenQueryCode() string {
	if filter.Cap() >= 1_000_000 {
		QueryCodeLen += 2
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		code := genCode(QueryCodeLen)
		if !filter.Test([]byte(code)) {
			filter.Add([]byte(code))
			return code
		}
	}

	for {
		code := genCode(QueryCodeLen + 2)
		if filter.Test([]byte(code)) {
			return code
		}
	}
}

func genCode(n int) string {
	var number string
	for i := 0; i < n; i++ {
		digit := rand.Intn(10)
		number += fmt.Sprintf("%d", digit)
	}
	return number
}

// func isExpired(code string) bool {
// 	var expireAt time.Time
// 	db.QueryRow("SELECT expire_at FROM countdown WHERE query_code = ?", code).Scan(&expireAt)
// 	return time.Now().After(expireAt)
// }

// 清理过期的查询码
// func codeCollect() {
// 	length := len(queryCodeSet)

// 	// 将queryCodeSet转换为slice
// 	codes := make([]string, 0, length)
// 	for code := range queryCodeSet {
// 		codes = append(codes, code)
// 	}

// 	// 设置随机种子
// 	rand.New(rand.NewSource(time.Now().UnixNano()))

// 	// 打乱slice的顺序
// 	rand.Shuffle(length, func(i, j int) {
// 		codes[i], codes[j] = codes[j], codes[i]
// 	})

// 	// 取前一部分元素进行回收操作
// 	var quarter int
// 	if length < CodeCollectThreshold>>4 {
// 		// 1.6w -> 8000
// 		quarter = length >> 1
// 	} else if length < CodeCollectThreshold>>8 {
// 		// 12.8w -> 8000
// 		quarter = length >> 4
// 	} else if length < CodeCollectThreshold>>10 {
// 		// 102.4w -> 8000
// 		quarter = length >> 8
// 	}
// 	for _, code := range codes[:quarter] {
// 		if isExpired(code) {
// 			queryCodeSet[code] = false
// 		}
// 	}
// }

// type Task func()

// 定时任务
// func startScheduledTask(task Task, period time.Duration, predicate func() bool) {
// 	ticker := time.NewTicker(period)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		if predicate() {
// 			go task()
// 		}
// 	}
// }
