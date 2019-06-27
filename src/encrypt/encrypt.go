package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
)

func padding(src []byte, blockSize int) ([]byte, int) {
	padNum := blockSize - len(src) % blockSize
	pad := bytes.Repeat([]byte{byte(0)}, padNum)
	return append(src, pad...), padNum
}

func unpadding(src []byte, padNum int) []byte {
	return src[:len(src) - padNum]
}

func EncryptAES(src []byte, key []byte) []byte {
	block,_ := aes.NewCipher(key)
	src, padNum := padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return append(src, byte(padNum))
}

func DecryptAES(src []byte, key[]byte) []byte {
	src, padNum := src[:len(src)-1], int(src[len(src)-1])
	block, _ := aes.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src[:len(src)-padNum]
}

func GetAESKey(key []byte) []byte {
	h := md5.Sum(key)
	return []byte(fmt.Sprintf("%x", h))[:16]
}
