package base256

import (
	"testing"
	"bytes"
)

func TestEncodingAndDecoding(t *testing.T) {
	src := []byte("Hello, World!")
	t.Logf("Original:\t%v", src)

	encoder := NewEncoding(StdEncoding)

	encodedData := encoder.EncodeToString(src)
	t.Logf("Encoded:\t%s", encodedData)

	encodedLen := encoder.EncodedLen(len(src))
	decodedLen := encoder.DecodedLen(encodedData)
	if encodedLen != decodedLen {
		t.Errorf("No Mached Length. encoded len: %d, decoded len: %d", encodedLen, decodedLen)
		t.FailNow()
	}

	decodedData := encoder.Decode(encodedData)
	t.Errorf("Decoded:\t%v", decodedData)

	if bytes.Compare(src, decodedData) != 0 {
		t.Errorf("No Matched Original Data.")
		t.FailNow()
	}
}