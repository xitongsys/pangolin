package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext) % blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func EncryptAES(origData, key []byte) ([]byte, error) {
	origData = []byte(base64.StdEncoding.EncodeToString(origData))
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    origData = PKCS7Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

func DecryptAES(crypted, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	origData, err = base64.StdEncoding.DecodeString(string(origData))
    return origData, err
}

func GetAESKey(key []byte) []byte {
	h := md5.Sum(key)
	return []byte(fmt.Sprintf("%x", h))[:16]
}
