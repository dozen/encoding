package base128

import (
	"unicode/utf8"
)

const (
	encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789あいうえおかきくけこさしすせそざじずぜぞたちつてとだぢづでどなにぬねのはひふへほばびぶべぼぱぴぷぺぽまみむめもやゆよらりるれろわをん"
	NormShift = 48
	EncShift  = 49
)

const StdEncoding = encodeStd

type Encoding struct {
	encode    []rune
	decodeMap map[int32]byte
}

func NewEncoding(encoder string) *Encoding {
	e := new(Encoding)
	e.encode = []rune(encoder)

	if len(e.encode) < 128 {
		panic("encoding is not 128-characters long")
	}

	i := 0
	e.decodeMap = map[int32]byte{}
	for _, v := range e.encode {
		e.decodeMap[v] = byte(i)
		i++
	}
	return e
}

func (e *Encoding) EncodedLen(n int) int {
	return n*8/7 + 1
}

func (e *Encoding) DecodedLen(src string) int {
	return utf8.RuneCountInString(src) * 7 / 8
}

func (e *Encoding) EncodeToString(src []byte) string {
	dst := make([]rune, e.EncodedLen(len(src)))
	encMap := e.encode

	di := 0
	si := 0
	n := (len(src) / 7) * 7
	for si < n {
		var val uint
		for i, sh := 0, uint(NormShift); i < 7; i++ {
			val |= uint(src[si+i]) << sh
			sh -= 8
		}

		for i, sh := 0, uint(EncShift); i < 8; i++ {
			dst[di+i] = encMap[val>>uint(sh)&0x7F]
			sh -= 7
		}
		si += 7
		di += 8
	}

	remain := len(src) - si
	if remain == 0 {
		return string(dst)
	}

	var val uint
	for i, sh := 0, uint(NormShift); i < remain; i++ {
		val |= uint(src[si+i]) << sh
		sh -= 8
	}

	for i, sh := 0, uint(EncShift); i <= remain; i++ {
		dst[di+i] = encMap[val>>uint(sh)&0x7F]
		sh -= 7
	}

	return string(dst)
}

func (e *Encoding) Decode(srcStr string) []byte {
	dstSize := e.DecodedLen(srcStr)
	//decMap := e.decodeMap
	var dst = make([]byte, dstSize)

	bufi := 0
	dstlen := 0
	var dbuf [8]byte
	for _, v := range srcStr {
		dbuf[bufi] = e.decodeMap[v]

		bufi++
		if bufi == 8 {
			var val uint
			for i, sh := 0, uint(EncShift); i < 8; i++ {
				val |= uint(dbuf[i]) << sh
				sh -= 7
			}

			for i, sh := 0, uint(NormShift); i < 7; i++ {
				dst[dstlen+i] = byte(val >> sh)
				sh -= 8
			}
			dstlen += 7
			bufi = 0
		}
	}

	if bufi == 0 {
		return dst
	}

	var val uint
	for i, sh := 0, uint(EncShift); i < 8; i++ {
		val |= uint(dbuf[i]) << sh
		sh -= 7
	}

	for i, sh := 0, uint(NormShift); dstlen+i < dstSize; i++ {
		dst[dstlen+i] = byte(val >> sh)
		sh -= 8
	}

	return dst
}
