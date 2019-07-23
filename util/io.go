package util

import (
	"io"
)

func ReadPacket(reader io.Reader) ([]byte, error) {
	data, lenBs := []byte{}, []byte{0}
	for {
		if _, err := ReadFull(reader, lenBs); err != nil {
			return nil, err
		}

		if ln := int(lenBs[0]); ln > 0 {
			cur := make([]byte, ln)
			if _, err := ReadFull(reader, cur); err != nil {
				return nil, err
			}
			data = append(data, cur...)

		} else {
			break
		}
	}
	return data, nil
}

func WritePacket(writer io.Writer, data []byte) (n int, err error) {
	n = len(data)
	for len(data) > 0 {
		wc := 255
		if len(data) < wc {
			wc = len(data)
		}

		if _, err := WriteFull(writer, []byte{byte(wc)}); err != nil {
			return n - len(data), err
		}
		if _, err := WriteFull(writer, data[:wc]); err != nil {
			return n - len(data), err
		}
		data = data[wc:]
	}
	_, err = WriteEnd(writer)
	return n - len(data), err
}

func ReadFull(reader io.Reader, buf []byte) (n int, err error) {
	ln, left := len(buf), len(buf)
	for left > 0 {
		if n, err = reader.Read(buf[ln-left:]); n > 0 && err == nil {
			left -= n
		} else if err != nil {
			break
		}
	}
	return ln - left, err
}

func WriteFull(writer io.Writer, buf []byte) (n int, err error) {
	ln, left := len(buf), len(buf)
	for left > 0 {
		if n, err = writer.Write(buf[ln-left:]); n > 0 && err == nil {
			left -= n
		} else if err != nil {
			break
		}
	}
	return ln - n, err
}
func WriteEnd(writer io.Writer) (n int, err error) {
	bs := []byte{0}
	return WriteFull(writer, bs)
}
