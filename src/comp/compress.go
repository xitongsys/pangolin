package comp

import (
	"bytes"
	"io/ioutil"
	"compress/gzip"
)

func reverse(buf []byte) []byte {
	n := len(buf)
	res := make([]byte, n)
	for i:=0; i<n; i++{
		res[n-1-i] = buf[i]
	}
	return res
}

func UncompressGzip(buf []byte) ([]byte, error) {
	rbuf := bytes.NewReader(buf)
	gzipReader, _ := gzip.NewReader(rbuf)
	res, err := ioutil.ReadAll(gzipReader)
	return res, err
}

func CompressGzip(buf []byte) []byte {
	buf = buf
	var res bytes.Buffer
	gzipWriter := gzip.NewWriter(&res)
	gzipWriter.Write(buf)
	gzipWriter.Close()
	return res.Bytes()
}
