package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabate = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomInt generates a random integer between min max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabate)
	for i := 0; i < n; i++ {
		c := alphabate[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//RandomOwner generates a random qwner name
func RandomOwner() string {
	return RandomString(6)
}

//RandomAccountno generates a random account number
func RandomAccountno() int64 {
	return RandomInt(123456789, 987456123987)
}

//RandomEmail generatesa random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
