package encrypt

import (
	"fmt"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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

func EncryptAES(origData, key []byte) (bs []byte, rerr error) {
	defer func() {
		if r := recover(); r!=nil {
			rerr = fmt.Errorf("Encrypt failed")
		}
	}()

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

func DecryptAES(crypted, key []byte) (bs []byte, rerr error) {
	defer func() {
		if r := recover(); r!=nil{
			rerr = fmt.Errorf("Decrypt failed")
		}
	}()

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
    return origData, err
}

func GetAESKey(key []byte) []byte {
	res := make([]byte, 16)
	for i:=0; i<16 && i<len(key); i++{
		res[i] = key[i]
	}
	for i:=len(key); i<16; i++ {
		res[i] = '0'
	}
	return res
}
