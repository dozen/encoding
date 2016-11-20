package base256

import (
	"unicode/utf8"
)

const (
	encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわをんァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヲンヴヵヶ亜唖娃阿哀愛挨姶逢葵茜穐悪握渥旭葦芦鯵梓圧斡扱宛姐虻飴絢綾"
	StdEncoding = encodeStd
)

type Encoding struct {
	encode    []rune
	decodeMap map[int32]byte
}

func NewEncoding(encoder string) *Encoding {
	e := new(Encoding)
	e.encode = []rune(encoder)

	if len(e.encode) < 256 {
		panic("encoding is not 256-characters long")
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
	return n
}

func (e *Encoding) DecodedLen(src string) int {
	return utf8.RuneCountInString(src)
}

func (e *Encoding) EncodeToString(src []byte) string {
	dst := make([]rune, e.EncodedLen(len(src)))
	encMap := e.encode

	for i, v := range src {
		dst[i] = encMap[v]
	}

	return string(dst)
}

func (e *Encoding) Decode(srcStr string) []byte {
	dstSize := e.DecodedLen(srcStr)
	var dst = make([]byte, dstSize)

	var i = 0
	for _, v := range srcStr {
		dst[i] = e.decodeMap[v]
		i++
	}

	return dst
}
