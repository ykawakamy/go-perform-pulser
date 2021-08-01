package util

import (
	"time"
)

// パフォーマンス測定用の構造体
type Ping struct {
	// シーケンス番号 topic単位に連番を設定する。
	seq uint32
	// 送信時のシステム時刻(ms)（正確にはこのメッセージの生成時刻）
	send_tick uint64
	// 受信時のシステム時刻(ms)（正確にはこのメッセージのデコード時刻）
	receive_tick uint64
	// receive_tick - send_tick
	latency uint64
}

func CreatePing(seq uint32) Ping {
	return Ping{
		seq:       seq,
		send_tick: getCurrentTick(),
	}
}

func (s Ping) Serialize(i int) []byte {
	// NOTE: 16進への変換はライブラリを使用すると実装の差が出ると考え、手書きした。
	buf := make([]byte, i)
	seq := int(s.seq)
	buf[0] = toHex(seq >> 28)
	buf[1] = toHex(seq >> 24)
	buf[2] = toHex(seq >> 20)
	buf[3] = toHex(seq >> 16)
	buf[4] = toHex(seq >> 12)
	buf[5] = toHex(seq >> 8)
	buf[6] = toHex(seq >> 4)
	buf[7] = toHex(seq >> 0)

	send_tick := int(s.send_tick)
	buf[8] = toHex(send_tick >> 60)
	buf[9] = toHex(send_tick >> 56)
	buf[10] = toHex(send_tick >> 52)
	buf[11] = toHex(send_tick >> 48)
	buf[12] = toHex(send_tick >> 44)
	buf[13] = toHex(send_tick >> 40)
	buf[14] = toHex(send_tick >> 36)
	buf[15] = toHex(send_tick >> 32)
	buf[16] = toHex(send_tick >> 28)
	buf[17] = toHex(send_tick >> 24)
	buf[18] = toHex(send_tick >> 20)
	buf[19] = toHex(send_tick >> 16)
	buf[20] = toHex(send_tick >> 12)
	buf[21] = toHex(send_tick >> 8)
	buf[22] = toHex(send_tick >> 4)
	buf[23] = toHex(send_tick >> 0)

	FillArray(buf[24:], 0x20)
	return buf
}

func DeserialPing(data []byte) Ping {
	tick := getCurrentTick()
	seq := HexToUInt(data[0:8])
	send_tick := HexToUInt(data[8:24])

	return Ping{
		seq:          uint32(seq),
		send_tick:    send_tick,
		receive_tick: tick,
		latency:      tick - send_tick,
	}
}

/* レイテンシ測定用にシステムのミリ秒を返す */
func getCurrentTick() uint64 {
	return uint64(time.Now().UnixNano() / 1000000) // milli sec.
}
