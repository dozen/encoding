package base128

import (
	"bytes"
	"testing"
)

func TestSet(t *testing.T) {
	str := []byte("T.I.A. This is Africa.")

	for i, _ := range str {
		EncodingAndDecoding(str[:i+1], t)
	}
}

func EncodingAndDecoding(src []byte, t *testing.T) {
	t.Logf("Original:\t%v", src)
	encoder := StdEncoding

	encodedData := encoder.EncodeToString(src)
	t.Logf("Encoded:\t%s", encodedData)

	decodedData, err := encoder.DecodeString(encodedData)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	t.Logf("Decoded:\t%v", decodedData)

	if bytes.Compare(src, decodedData) != 0 {
		t.Errorf("No Matched Original Data.")
		t.FailNow()
	}
}
