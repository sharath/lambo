package authentication

import (
	"encoding/base64"
)

// CookieCoder is a base64 encoding that a web browser can store as a cookie
var CookieCoder = base64.URLEncoding.WithPadding('*')

// Encode CookieEncodes an input byte array and returns the result
func Encode(input []byte) []byte {
	var output []byte
	output = []byte(CookieCoder.EncodeToString(input))
	return output
}

// Decode Decodes an input byte array and returns the result
func Decode(input []byte) []byte {
	var output []byte
	output, _ = CookieCoder.DecodeString(string(input[:]))
	return output
}
