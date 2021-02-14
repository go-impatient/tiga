package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

// 解密, key 长度必须为16 byte, js是
func AesDecrypt(hexS, key []byte) ([]byte, error) {
	// 转成字节数组
	hexRaw, err := hex.DecodeString(string(hexS))
	if err != nil {
		return nil, err
	}
	if len(key) == 0 {
		return nil, errors.New("key 不能为空")
	}
	pkey := paddingLeft(key, '0', 16)
	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return nil, fmt.Errorf("key 长度必须 16/24/32长度: %s", err)
	}
	// 加密模式
	blockModel := cipher.NewCBCDecrypter(block, pkey)
	plantText := make([]byte, len(hexRaw))
	// 解密
	blockModel.CryptBlocks(plantText, hexRaw)
	// 去补全码
	plantText = pkcs7UnPadding(plantText)
	return plantText, nil
}

// 加密
func AesEncrypt(raw string, key string) (string, error) {
	// 转换成字节数组
	origData := []byte(raw)
	if len(key) == 0 {
		return "", errors.New("key 不能为空")
	}
	k := paddingLeft([]byte(key), '0', 16)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", fmt.Errorf("填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256  key 长度必须 16/24/32长度: %s", err)
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k)
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中会导致错误
	return hex.EncodeToString(cryted), nil
}

func paddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}

// pkcs7Padding 补码
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs7UnPadding 去码
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
