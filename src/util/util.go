package util

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

const (
	key = "uPpqu1JEjhHUzgKxAljqfY6RZEIiEJEe568QUJEFbIkyWIUz1lKqvP1e21pXjRk1"
)

// CalculateHash calculate hash
func CalculateHash(data string) string {
	h := sha256.New()
	h.Write([]byte(key + data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// RandomString return a random string with length len
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
