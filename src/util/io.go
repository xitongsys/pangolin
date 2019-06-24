package util

import (
	"io"
)

func ReadPacket(reader io.Reader) []byte {
	data, lenBs := []byte{}, []byte{0}
	for {
		ReadFull(reader, lenBs)
		if ln := int(lenBs[0]); ln > 0 {
			cur := make([]byte, ln)
			ReadFull(reader, cur)
			data = append(data, cur...)

		}else{
			break
		}
	}
	return data
}

func WritePacket(writer io.Writer, data []byte) {
	for len(data) > 0 {
		wc := 255
		if len(data) < wc {
			wc = len(data)
		}

		WriteFull(writer, []byte{byte(wc)})
		WriteFull(writer, data[:wc])
		data = data[wc:]
	}
	WriteEnd(writer)
}

func ReadFull(reader io.Reader, buf []byte){
	ln, left := len(buf), len(buf)
	for left > 0 {
		if n, err := reader.Read(buf[ln - left:]); n > 0 && err == nil {
			left -= n
		}
	}
}

func WriteFull(writer io.Writer, buf []byte){
	ln, left := len(buf), len(buf)
	for left > 0 {
		if n, err := writer.Write(buf[ln - left:]); n > 0 && err == nil {
			left -= n
		}
	}
}
func WriteEnd(writer io.Writer){
	bs := []byte{0}
	WriteFull(writer, bs)
}