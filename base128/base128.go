package base128

import (
	"io"
	"unicode/utf8"
)

const (
	encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789あいうえおかきくけこさしすせそざじずぜぞたちつてとだぢづでどなにぬねのはひふへほばびぶべぼぱぴぷぺぽまみむめもやゆよらりるれろわをん"
	normShift = 48
	encShift  = 49
)

var StdEncoding = NewEncoding(encodeStd)

type Encoding struct {
	encode    [128]rune
	decodeMap map[int32]byte
}

func NewEncoding(encoder string) *Encoding {
	if utf8.RuneCountInString(encoder) < 128 {
		panic("encoding is not 128-characters long")
	}

	e := new(Encoding)
	copy(e.encode[:], []rune(encoder))

	i := 0
	e.decodeMap = map[int32]byte{}
	for _, v := range e.encode {
		e.decodeMap[v] = byte(i)
		i++
	}
	return e
}

func (e *Encoding) EncodedLen(n int) int {
	return (n*8 + 6) / 7
}

func (e *Encoding) DecodedLen(src string) int {
	return utf8.RuneCountInString(src) * 7 / 8
}

func (e *Encoding) Encode(dst []rune, src []byte) {
	if len(src) == 0 {
		return
	}

	encMap := e.encode

	di := 0
	si := 0
	n := (len(src) / 7) * 7
	for si < n {
		var val uint64
		for i, sh := 0, uint(normShift); i < 7; i++ {
			val |= uint64(src[si+i]) << sh
			sh -= 8
		}

		for i, sh := 0, uint(encShift); i < 8; i++ {
			dst[di+i] = encMap[val>>uint(sh)&0x7F]
			sh -= 7
		}
		si += 7
		di += 8
	}

	remain := len(src) - si
	if remain == 0 {
		return
	}

	var val uint64
	for i, sh := 0, uint(normShift); i < remain; i++ {
		val |= uint64(src[si+i]) << sh
		sh -= 8
	}

	for i, sh := 0, uint(encShift); i <= remain; i++ {
		dst[di+i] = encMap[val>>uint(sh)&0x7F]
		sh -= 7
	}
}

func (e *Encoding) EncodeToString(src []byte) string {
	dst := make([]rune, e.EncodedLen(len(src)))
	e.Encode(dst, src)
	return string(dst)
}

func (e *Encoding) DecodeString(s string) ([]byte, error) {
	dst := make([]byte, e.DecodedLen(s))
	e.decode(dst, []rune(s))
	return dst, nil
}

func (e *Encoding) Decode(dst []byte, src []rune) int {
	return e.decode(dst, src)
}

func (e *Encoding) decode(dst []byte, src []rune) (n int) {
	n = len(dst)
	bufi := 0
	dstlen := 0
	var dbuf [8]byte
	for _, v := range src {
		dbuf[bufi] = e.decodeMap[v]

		bufi++
		if bufi == 8 {
			var val uint64
			for i, sh := 0, uint(encShift); i < 8; i++ {
				val |= uint64(dbuf[i]) << sh
				sh -= 7
			}

			for i, sh := 0, uint(normShift); i < 7; i++ {
				dst[dstlen+i] = byte(val >> sh)
				sh -= 8
			}
			dstlen += 7
			bufi = 0
		}
	}

	if bufi == 0 {
		return
	}

	var val uint64
	for i, sh := 0, uint(encShift); i < 8; i++ {
		val |= uint64(dbuf[i]) << sh
		sh -= 7
	}

	for i, sh := 0, uint(normShift); dstlen+i < n; i++ {
		dst[dstlen+i] = byte(val >> sh)
		sh -= 8
	}

	return
}

type encoder struct {
	err  error
	enc  *Encoding
	w    io.Writer
	buf  [7]byte
	nbuf int
	out  [1024]rune
}

func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	return &encoder{enc: enc, w: w}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	if e.nbuf > 0 {
		var i int
		for i = 0; i < len(p) && e.nbuf < 7; i++ {
			e.buf[e.nbuf] = p[i]
			e.nbuf++
		}
		n += i
		p = p[i:]
		if e.nbuf < 7 {
			return
		}
		e.enc.Encode(e.out[:], e.buf[:])

		if _, e.err = e.w.Write([]byte(string(e.out[:8]))); e.err != nil {
			return n, e.err
		}
		e.nbuf = 0
	}

	for len(p) >= 7 {
		nn := len(e.out) / 8 * 7
		if nn > len(p) {
			nn = len(p)
			nn -= nn % 7
		}
		e.enc.Encode(e.out[:], p[:nn])
		if _, e.err = e.w.Write([]byte(string(e.out[0 : nn/7*8]))); e.err != nil {
			return n, e.err
		}
		n += nn
		p = p[nn:]
	}

	for i := 0; i < len(p); i++ {
		e.buf[i] = p[i]
	}
	e.nbuf = len(p)
	n += len(p)
	return
}

func (e *encoder) Close() error {
	if e.err == nil && e.nbuf > 0 {
		e.enc.Encode(e.out[:], e.buf[:e.nbuf])
		_, e.err = e.w.Write([]byte(string(e.out[:e.enc.EncodedLen(e.nbuf)])))
		e.nbuf = 0
	}
	return e.err
}
