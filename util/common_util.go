package util

// math.max()のuint32版 ※golangの標準ライブラリにないため実装
func max(x, y uint32) uint32 {
	if x > y {
		return x
	} else {
		return y
	}

}

// math.max()のuint32版 ※golangの標準ライブラリにないため実装
func min(x, y uint32) uint32 {
	if x < y {
		return x
	} else {
		return y
	}
}

// 入力した数値の下位4bitを16進(ASCIIコードの0x30-0x39, 0x61-0x66)に変換する。
func toHex(v int) byte {
	h := v & 0x0f
	if h < 10 {
		return byte(0x30 + h)
	} else {
		return byte(0x61 + h - 10)
	}

}

// 16進(ASCIIコードの0x30-0x39, 0x61-0x66)を4bit数値に変換する
func toDec4bit(v byte) int64 {
	if 0x30 <= v && v <= 0x39 {
		return int64(v - 0x30)
	}
	if 0x61 <= v && v <= 0x66 {
		return int64(v - 0x61 + 10)
	}
	panic("invalid value.")

}

// バイト列をint64に変換します。
// アライメントに合わせて入力バイト列はスライスしてください。
// ex) int32ならb[0:4]
func HexToInt(buf []byte) int64 {
	var h int64 = 0
	for _, it := range buf {
		h = h<<4 | toDec4bit(it)
	}
	return h
}

// バイト列をuint64に変換します。
// アライメントに合わせて入力バイト列はスライスしてください。
// ex) int32ならb[0:4]
func HexToUInt(buf []byte) uint64 {
	var h uint64 = 0
	for _, it := range buf {
		h = h<<4 | uint64(toDec4bit(it))
	}
	return h
}

// バイト列を指定した値で埋めます。
func FillArray(buf []byte, padder byte) {
	if len(buf) == 0 {
		return
	}
	// @see https://gist.github.com/taylorza/df2f89d5f9ab3ffd06865062a4cf015d
	buf[0] = padder

	for j := 1; j < len(buf); j *= 2 {
		copy(buf[j:], buf[:j])
	}
}
