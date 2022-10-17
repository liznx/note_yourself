package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// 算法名
const KeyName string = "AES"

// 加解密算法/模式/填充方式
// ECB模式只用密钥即可对数据进行加密解密，CBC模式需要添加一个iv
const Cipher_algorithma string = "AES/CBC/PKCS7Padding"

/**
* 微信 数据解密<br/>
* 对称解密使用的算法为 AES-128-CBC，数据采用PKCS#7填充<br/>
* 对称解密的目标密文:encrypted=Base64_Decode(encryptData)<br/>
* 对称解密秘钥:key = Base64_Decode(session_key),aeskey是16字节<br/>
* 对称解密算法初始向量:iv = Base64_Decode(iv),同样是16字节<br/>
*
* @param encrypted 目标密文
* @param session_key 会话ID
* @param iv 加密算法的初始向量
 */

func DecryptWXOpenData(app_id string, sessionKey string, encryptData string, iv string) (map[string]interface{}, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	dataBytes, _ := AesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)
	fmt.Println(string(dataBytes))
	m := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &m)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	temp := m["watermark"].(map[string]interface{})
	appid := temp["appid"].(string)
	if appid != app_id {
		return nil, fmt.Errorf("invalid appid, get !%s!", appid)
	}
	if err != nil {
		return nil, err
	}
	return m, nil

}

func AesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	//获取的数据尾端有'/x0e'占位符,去除它
	for i, ch := range origData {
		if ch == '\b' {
			origData[i] = ' '
		}
	}
	return origData, nil
}
