package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().Unix())
}

// generate a random integer between min and max
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		pos := rand.Int() % k
		sb.WriteByte(alphabet[pos])
	}

	return sb.String()
}

// generate a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "NTD"}
	n := len(currencies)
	return currencies[rand.Int()%n]
}

// generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
