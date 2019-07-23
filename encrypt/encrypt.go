package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func EncryptAES(origData []byte, keyBytes []byte) (bs []byte, rerr error) {
	defer func() {
		if r := recover(); r != nil {
			rerr = fmt.Errorf("Eecrypt failed")
		}
	}()

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DecryptAES(crypted []byte, keyBytes []byte) (bs []byte, rerr error) {
	defer func() {
		if r := recover(); r != nil {
			rerr = fmt.Errorf("Decrypt failed")
		}
	}()

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func GetAESKey(key []byte) []byte {
	res := make([]byte, 16)
	for i := 0; i < 16 && i < len(key); i++ {
		res[i] = key[i]
	}
	return res
}
