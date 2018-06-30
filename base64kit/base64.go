package base64kit

import (
	"encoding/base64"
)

const (
	temp string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var base64Temp *base64.Encoding = nil

func NewEncoding() *base64.Encoding {
	if base64Temp == nil {
		base64Temp = base64.NewEncoding(temp)
	}

	return base64Temp
}

func Base64EncodeByByte(src []byte) string {
	return NewEncoding().EncodeToString(src)
}

func Base64EncodeByString(str string) string {
	return Base64EncodeByByte([]byte(str))
}

func Base64Decode(str string) (string, error) {
	b, err := NewEncoding().DecodeString(str)
	if len(b) < 1 {
		return "", err
	} else {
		return string(b), nil
	}

}

