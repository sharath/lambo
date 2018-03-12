package authentication

import (
	"encoding/base64"
)

var CookieCoder = base64.URLEncoding.WithPadding('*')

func Encode(input []byte) []byte {
	var output []byte
	output = []byte(CookieCoder.EncodeToString(input))
	return output
}

func Decode(input []byte) []byte {
	var output []byte
	output, _ = CookieCoder.DecodeString(string(input[:]))
	return output
}
