package msggenerator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers       = "0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	tlds = []string{"com", "org", "net", "io", "dev"}
)

func randInt(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
}

func randString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx := randInt(int64(len(letters)))
		sb.WriteByte(letters[idx])
	}
	return sb.String()
}

func randProto() string {
	if randInt(2) == 0 {
		return "http"
	}
	return "https"
}

// GenerateRandomURL generates a random URL with random subdomain, domain, TLD, path and query parameters
func GenerateRandomURL() string {
	protocol := randProto()
	subdomain := randString(5)
	domain := randString(8)
	tld := tlds[randInt(int64(len(tlds)))]
	path := "/" + randString(10)
	query := "?" + randString(5) + "=" + randString(5)

	return fmt.Sprintf("%s://%s.%s.%s%s%s", protocol, subdomain, domain, tld, path, query)
}
