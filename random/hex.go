package random

import (
	rand "crypto/rand"
	"encoding/hex"
)

//Hex creates l long random string encoded in hex. Uses crypto/rand
func Hex(l int64) string {
	b := make([]byte, l)
	for {
		n, err := rand.Read(b[:cap(b)])
		if err != nil {
			continue
		}
		b = b[:n]
		if len(b) != int(l) {
			continue
		}
		break
	}

	return hex.EncodeToString(b)[:l]
}
