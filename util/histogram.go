package util

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

type Histogram struct {
	buckets *sync.Map
	m       *sync.Mutex
}

func CreateHistogram() *Histogram {
	return &Histogram{
		buckets: new(sync.Map),
		m:       new(sync.Mutex),
	}
}

func (s *Histogram) IncreamentPing(ping *Ping) {
	s.Increament(ping.latency)
}

func (s *Histogram) Increament(v uint64) {
	target := convertBucketKey(v)
	val, _ := s.buckets.Load(target)

	var pc *uint32
	if val == nil {
		s.m.Lock()
		val, _ := s.buckets.Load(target)
		if val == nil {
			pc = new(uint32)
			*pc = 0
			s.buckets.Store(target, pc)
		} else {
			pc = val.(*uint32)
		}
		s.m.Unlock()
	} else {
		pc = val.(*uint32)
	}

	atomic.AddUint32(pc, 1)
}

func (s *Histogram) Print() {
	fmt.Printf("--- latancy histogram ---\n")
	s.buckets.Range(func(k, v interface{}) bool {
		mn := deconvertBucketKey(k.(int))
		mx := deconvertBucketKey(k.(int)+1) - 1
		fmt.Printf("%d - %d : %d\n", mn, mx, *v.(*uint32))
		return true
	})

}

// 数値からヒストグラムのBuckectのインデックスを計算する。
//
// TODO: 現状、固定で対数ヒストグラムとしているが、ストラテジパターンで変更できるようにしたい。
func convertBucketKey(v uint64) int {
	return int(math.Log(float64(v)))
}

// ヒストグラムのBuckectのインデックスから推定される最小値を計算する。
//
// NOTE: 推定される最大値を取得する場合、`deconvertBucketKey(i+1) - 1`とする。
func deconvertBucketKey(u int) int {
	return int(math.Pow(math.E, float64(u)))
}
