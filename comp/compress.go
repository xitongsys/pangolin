package comp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

func reverse(buf []byte) []byte {
	n := len(buf)
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[n-1-i] = buf[i]
	}
	return res
}

func UncompressGzip(buf []byte) (bs []byte, rerr error) {
	return bs, nil
	defer func() {
		if err := recover(); err != nil {
			rerr = fmt.Errorf("%v", err)
		}
	}()
	rbuf := bytes.NewReader(buf)
	gzipReader, _ := gzip.NewReader(rbuf)
	res, err := ioutil.ReadAll(gzipReader)
	return res, err
}

func CompressGzip(buf []byte) []byte {
	return buf
	buf = buf
	var res bytes.Buffer
	gzipWriter := gzip.NewWriter(&res)
	gzipWriter.Write(buf)
	gzipWriter.Close()
	return res.Bytes()
}
