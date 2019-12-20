package utils

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func ZlibMarshal(p []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(p)
	w.Close()
	return b.Bytes(), nil
}

func ZlibUnmarshal(p []byte) ([]byte, error) {
	var b bytes.Buffer
	b.Write(p)
	r, e := zlib.NewReader(&b)
	if e != nil {
		return nil, e
	}

	return ioutil.ReadAll(r)
}
