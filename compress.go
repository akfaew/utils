package utils

import (
	"bytes"
	"compress/gzip"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

func Compress(data interface{}) []byte {
	m, err := msgpack.Marshal(data)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(m); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	return b.Bytes()
}

func Decompress(enc []byte, out interface{}) {
	gr, err := gzip.NewReader(bytes.NewReader(enc))
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	if err := msgpack.NewDecoder(gr).Decode(out); err != nil {
		panic(err)
	}
}
