package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var (
	M3u8Key = []byte("ABCDEFGHIJKLMNOP")
	Key     = []byte("HEoq8f4N8lLriYIA")
	Auth    = []byte("5f570a7d150f60bI")
)

func AesEncryptCBC(origData []byte, key []byte) string {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted := make([]byte, len(origData))                    // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return hex.EncodeToString(encrypted)
}

func AesDecryptCBC(str string, key []byte) (err error, data []byte) {
	encrypted, err := hex.DecodeString(str)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted := make([]byte, len(encrypted))                   // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	data = pkcs5UnPadding(decrypted)                            // 去除补全码
	return
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
