package util

import (
	"fmt"
	"io"
	"net"
	"time"
)

func WriteUntil(conn net.Conn, bufSize int, data []byte, timeout time.Duration, fun func([]byte) bool) (res []byte, err error) {
	after := time.After(timeout)
	done := make(chan bool)
	buf := make([]byte, bufSize)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
		close(done)
	}()

	go func() {
		for {
			_, errW := conn.Write(data)
			if errW != nil {
				err = errW
				return
			}

			select {
			case <-done:
				return
			case <-after:
				err = fmt.Errorf("timeout")
				conn.Close()
				return
			default:
			}

			time.Sleep(time.Millisecond * 100)
		}
	}()

	for err == nil {
		n, errR := conn.Read(buf)
		if errR != nil {
			err = errR
			return nil, err
		}

		if n <= 0 {
			continue
		}

		if fun(buf[:n]) {
			return buf[:n], nil
		}
	}
	return nil, err
}

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
