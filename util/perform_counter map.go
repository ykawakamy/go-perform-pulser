package util

import (
	"fmt"
	"sync"
)

type PerformCounterMap struct {
	m *sync.Mutex

	syncMap *sync.Map
}

func CreatePerformCounterMap() *PerformCounterMap {
	return &PerformCounterMap{m: new(sync.Mutex), syncMap: &sync.Map{}}
}

func (s *PerformCounterMap) Perform(topic string, ping *Ping) {
	val, _ := s.syncMap.Load(topic)

	var pc *PerformCounter
	if val == nil {
		s.m.Lock()
		val, _ := s.syncMap.Load(topic)
		if val == nil {
			pc = CreatePerformCounter()
			s.syncMap.Store(topic, pc)
		} else {
			pc = val.(*PerformCounter)
		}
		s.m.Unlock()
	} else {
		pc = val.(*PerformCounter)
	}
	pc.Perform(ping.seq, ping.latency)
}

func (s *PerformCounterMap) CollectAndReset() PerformSnapShot {
	snapSum := CreatePerformSummary()

	s.syncMap.Range(func(k, v interface{}) bool {
		pc := v.(*PerformCounter)
		snap := pc.CollectAndReset()
		snapSum.add(&snap)
		return true
	})

	return snapSum
}

func (s *PerformCounterMap) Size() int {
	cnt := 0
	s.syncMap.Range(func(k, v interface{}) bool {
		cnt++
		fmt.Printf("%v %v\n", k, v)
		return true
	})

	return cnt
}
