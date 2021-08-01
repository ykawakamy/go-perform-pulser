package util

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PerformNormal(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(0, 0)
	pc.Perform(1, 0)
	pc.Perform(2, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(3), pc.expected_seq)
	assert.Equal(t, uint(0), snap.no_update)
	assert.Equal(t, uint(0), snap.rewind_error)
	assert.Equal(t, uint(0), snap.skip_error)
}

func Test_PerformNoUpdate(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(0, 0)
	pc.Perform(0, 0)
	pc.Perform(0, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(1), pc.expected_seq)
	assert.Equal(t, uint(2), snap.no_update)
	assert.Equal(t, uint(0), snap.rewind_error)
	assert.Equal(t, uint(0), snap.skip_error)
}

func Test_PerformSkip(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(1, 0)
	pc.Perform(1, 0)
	pc.Perform(1, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(2), pc.expected_seq)
	assert.Equal(t, uint(2), snap.no_update)
	assert.Equal(t, uint(0), snap.rewind_error)
	assert.Equal(t, uint(1), snap.skip_error)
}

func Test_PerformRewind_Sequencial(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(2, 0)
	pc.Perform(1, 0)
	pc.Perform(0, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(1), pc.expected_seq)
	assert.Equal(t, uint(0), snap.no_update)
	assert.Equal(t, uint(2), snap.rewind_error)
	assert.Equal(t, uint(1), snap.skip_error)
}

func Test_PerformError_SNR(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(2, 0)
	pc.Perform(2, 0)
	pc.Perform(0, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(1), pc.expected_seq)
	assert.Equal(t, uint(1), snap.no_update)
	assert.Equal(t, uint(1), snap.rewind_error)
	assert.Equal(t, uint(1), snap.skip_error)
}

func Test_PerformError_SNO(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(2, 0)
	pc.Perform(0, 0)
	pc.Perform(1, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(2), pc.expected_seq)
	assert.Equal(t, uint(0), snap.no_update)
	assert.Equal(t, uint(1), snap.rewind_error)
	assert.Equal(t, uint(1), snap.skip_error)
}

func Test_PerformError_SRS(t *testing.T) {
	pc := CreatePerformCounter()

	pc.Perform(2, 0)
	pc.Perform(0, 0)
	pc.Perform(2, 0)

	snap := pc.CollectAndReset()
	assert.Equal(t, uint32(3), pc.expected_seq)
	assert.Equal(t, uint(0), snap.no_update)
	assert.Equal(t, uint(1), snap.rewind_error)
	assert.Equal(t, uint(2), snap.skip_error)
}

func Test_PerformMap(t *testing.T) {
	pc := CreatePerformCounterMap()

	wg := new(sync.WaitGroup)
	thread := 100
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		key := strconv.Itoa(i)
		go func() {
			for j := 0; j < 1000; j++ {
				pc.Perform(key, &Ping{
					seq: uint32(j),
				})
			}
			wg.Done()
		}()
	}

	wg.Wait()

	cnt := 0
	pc.syncMap.Range(func(k, v interface{}) bool {
		cnt++
		return true
	})
	assert.Equal(t, thread, cnt)

	snap := pc.CollectAndReset()
	snap.Print(100)
}
