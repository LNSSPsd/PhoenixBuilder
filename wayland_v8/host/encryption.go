package host

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	//"rogchap.com/v8go"
)

//
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


//AesEncrypt 加密函数
func aesEncrypt(_plaintext, _key string) (string,string, error) {
	plaintext:=[]byte(_plaintext)
	key := []byte(_key)
	key32:=make([]byte,32)
	copy(key32,key)
	c := make([]byte, aes.BlockSize+len(plaintext))
	iv := c[:aes.BlockSize]

	block, err := aes.NewCipher(key32)
	if err != nil {
		return "","", err
	}
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	return hex.EncodeToString(crypted),hex.EncodeToString(iv), nil
}

// AesDecrypt 解密函数
func aesDecrypt(_ciphertext, _key, _iv string) (string, error) {
	ciphertext, _ :=hex.DecodeString(_ciphertext)
	key := []byte(_key)
	iv, _ :=hex.DecodeString(_iv)
	key32:=make([]byte,32)
	copy(key32,key)
	block, err := aes.NewCipher(key32)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = PKCS7UnPadding(origData)
	return string(origData), nil
}