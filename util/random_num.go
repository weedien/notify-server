package util

import (
	"math/rand"
)

// 生成一个随机数,该随机数在指定的范围内,并且在每个子区间的概率不同
func RandomRange(start, end, numRanges int, probabilityRatio float64) int {
	// 计算每个子区间的长度
	rangeLength := (end - start + 1) / numRanges

	// 计算每个子区间的概率
	probabilities := make([]float64, numRanges)
	probabilities[0] = 1.0
	for i := 1; i < numRanges; i++ {
		probabilities[i] = probabilities[i-1] * probabilityRatio
	}
	totalProbability := 0.0
	for _, p := range probabilities {
		totalProbability += p
	}
	for i := range probabilities {
		probabilities[i] /= totalProbability
	}

	// 生成随机数
	randomValue := rand.Float64()
	cumulativeProbability := 0.0
	for i, p := range probabilities {
		cumulativeProbability += p
		if randomValue <= cumulativeProbability {
			return start + i*rangeLength + rand.Intn(rangeLength)
		}
	}

	// 这种情况不应该发生,但如果发生,返回结束值
	return end
}
