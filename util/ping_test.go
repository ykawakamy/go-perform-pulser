package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Serialize(t *testing.T) {
	ping := CreatePing(1)

	buf := ping.Serialize(100)
	time.Sleep(10 * time.Millisecond)
	deping := DeserialPing(buf)

	assert.Equal(t, ping.seq, deping.seq)
	assert.Equal(t, ping.send_tick, deping.send_tick)
	assert.NotEqual(t, ping.receive_tick, deping.receive_tick)
	assert.NotEqual(t, ping.latency, deping.latency)

}
