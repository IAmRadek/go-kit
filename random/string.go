package random

import (
	"crypto/rand"
)

type Chars int

const (
	AlphaLarge Chars = 1 << iota
	AlphaSmall
	ZeroNine
)

func String(length int, chars Chars) string {
	if length <= 0 {
		return ""
	}

	var charset string
	if chars&AlphaLarge != 0 {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if chars&AlphaSmall != 0 {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}
	if chars&ZeroNine != 0 {
		charset += "0123456789"
	}

	if len(charset) == 0 {
		return ""
	}

	result := make([]byte, length)
	buf := make([]byte, length)

	for {
		if _, err := rand.Read(buf); err != nil {
			continue
		}

		for i := 0; i < length; i++ {
			result[i] = charset[int(buf[i])%len(charset)]
		}
		break
	}

	return string(result)
}
