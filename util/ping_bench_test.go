package util

import (
	"fmt"
	"math"
	"testing"
)

const (
	msgSize = 150
)

func Benchmark_SprintfAndAppendPadding(b *testing.B) {
	var cnt int = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str := fmt.Sprintf("%08x%016x", i, i)
		buf := []byte(str)
		padding := make([]byte, msgSize-len(buf))
		FillArray(buf, 0x20)
		buf = append(buf, padding...)
		cnt += len(buf)
	}
}

func Benchmark_AllocateOnce(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, msgSize)
		buf[0] = toHex(i >> 28)
		buf[1] = toHex(i >> 24)
		buf[2] = toHex(i >> 20)
		buf[3] = toHex(i >> 16)
		buf[4] = toHex(i >> 12)
		buf[5] = toHex(i >> 8)
		buf[6] = toHex(i >> 4)
		buf[7] = toHex(i >> 0)

		buf[8] = toHex(i >> 60)
		buf[9] = toHex(i >> 56)
		buf[10] = toHex(i >> 52)
		buf[11] = toHex(i >> 48)
		buf[12] = toHex(i >> 44)
		buf[13] = toHex(i >> 40)
		buf[14] = toHex(i >> 36)
		buf[15] = toHex(i >> 32)
		buf[16] = toHex(i >> 28)
		buf[17] = toHex(i >> 24)
		buf[18] = toHex(i >> 20)
		buf[19] = toHex(i >> 16)
		buf[20] = toHex(i >> 12)
		buf[21] = toHex(i >> 8)
		buf[22] = toHex(i >> 4)
		buf[23] = toHex(i >> 0)
	}
}

func Benchmark_log(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math.Logb(float64(i*10 + 1))
	}
}
