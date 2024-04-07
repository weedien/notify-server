package test

import (
	"testing"

	"github.com/weedien/countdown-server/infra/util"
)

func BenchmarkRandomNum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := util.RandomRange(7, 35, 3, 0.5)
		b.Logf("Random number: %d\n", num)
	}
}
