package utils

import "encoding/base64"

func EncodeToBase64(toEncode []byte) string {
	return base64.RawURLEncoding.EncodeToString(toEncode)
}

func DecodeString(toDecode string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(toDecode)
}
