package keygen

import "math/rand"

type KeyGenerator struct {
}

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func (k KeyGenerator) Generate(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(chars) {
			b[i] = chars[idx]
			i++
		}
	}
	return string(b)
}
