package base128

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncoder(t *testing.T) {
	str := "T.I.A. This is Africa."

	buf := bytes.NewBuffer([]byte{})
	w := NewEncoder(StdEncoding, buf)
	w.Write([]byte(str))
	w.Close()

	enc := string(buf.Bytes())
	fmt.Println(enc, "\n"+StdEncoding.EncodeToString([]byte(str)))

	dec, _ := StdEncoding.DecodeString(enc)
	fmt.Println(string(dec))
}

func TestEncodingAndDecoding(t *testing.T) {
	test := func(src []byte) {
		t.Logf("Original:\t%v", src)
		encoder := StdEncoding

		encodedData := encoder.EncodeToString(src)

		decodedData, err := encoder.DecodeString(encodedData)
		t.Logf("Decoded:\t%v", decodedData)
		t.Logf("Encoded:\t%s", encodedData)
		if err != nil {
			t.Errorf(err.Error())
			t.FailNow()
		}
		if bytes.Compare(src, decodedData) != 0 {
			t.Errorf("No Matched Original Data.")
			t.FailNow()
		}

	}

	src := "T.I.A. This is Africa."
	for i, _ := range src {
		test([]byte(src[i:]))
	}
}
