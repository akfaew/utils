package utils

import (
	"bytes"
	"compress/gzip"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

func Compress(data any) []byte {
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

func Decompress(enc []byte, out any) {
	if len(enc) == 0 {
		return
	}

	gr, err := gzip.NewReader(bytes.NewReader(enc))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := gr.Close(); err != nil {
			panic(err)
		}
	}()

	if err := msgpack.NewDecoder(gr).Decode(out); err != nil {
		panic(err)
	}
}
