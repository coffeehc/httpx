// cryto
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func CheckSum(msg []byte) uint16 {
	sum := 0
	for n := 0; n < len(msg); n++ {
		sum += int(msg[n])
	}
	//sum = (sum >> 16) + (sum & 0xffff)
	//sum += (sum >> 16)
	return uint16(^sum)
}

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

var aes_cfbiv = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}

func AesEncrypt(srctext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, aes_cfbiv)
	ciphertext := make([]byte, len(srctext))
	cfb.XORKeyStream(ciphertext, srctext)
	return ciphertext, nil
}

func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBDecrypter(block, aes_cfbiv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
