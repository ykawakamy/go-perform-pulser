package util

import (
	"math"
	"sync"
)

const (
	const_init_max = 0
	const_init_min = math.MaxUint32
)

type PerformCounter struct {
	expected_seq uint32
	perform      uint
	no_update    uint
	skip_error   uint
	rewind_error uint
	latency      uint64

	actual_max uint32
	actual_min uint32

	m *sync.Mutex
}

func CreatePerformCounter() *PerformCounter {
	return &PerformCounter{
		actual_max: const_init_max,
		actual_min: const_init_min,

		m: new(sync.Mutex),
	}
}

func (s *PerformCounter) Perform(actual uint32, latency uint64) {
	s.m.Lock()
	if s.expected_seq == actual {
		// success
	} else if s.expected_seq == actual+1 {
		s.no_update++
	} else if s.expected_seq < actual {
		s.skip_error++
	} else {
		s.rewind_error++
	}
	s.perform++

	s.expected_seq = actual + 1

	s.actual_max = max(s.actual_max, actual)
	s.actual_min = min(s.actual_min, actual)

	s.latency += latency

	s.m.Unlock()
}

func (s *PerformCounter) CollectAndReset() PerformSnapShot {
	s.m.Lock()
	snap := CreatePerformSnapShot(s)

	s.perform = 0
	s.no_update = 0
	s.skip_error = 0
	s.rewind_error = 0
	s.latency = 0

	s.actual_max = const_init_max
	s.actual_min = const_init_min

	s.m.Unlock()

	return snap
}
