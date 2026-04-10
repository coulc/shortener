package utils

import (
	"crypto/md5"
	"encoding/binary"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateShortCode(url string) string {
	hash := md5.Sum([]byte(url)) 
	num  := binary.BigEndian.Uint32(hash[:4])
	return encodeBase64(num)
}

func encodeBase64(num uint32) string {
	if num == 0 {
		return "0"
	}	

	base := uint32(len(charset))

var result []byte
	for num > 0 {
		remainder := num % base
		result = append([]byte{charset[remainder]},result...)
		num = num / base
	}
	return string(result)
}
