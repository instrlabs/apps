package models

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// GenerateUniqueUsername generates a unique username from email
func GenerateUniqueUsername(email string) (string, error) {
	base := email
	if at := strings.Index(email, "@"); at != -1 {
		base = email[:at]
	}
	base = strings.ToLower(strings.TrimSpace(base))
	if base == "" {
		base = "user"
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	suffix := fmt.Sprintf("%04d", nBig.Int64())
	candidate := fmt.Sprintf("%s%s", base, suffix)

	return candidate, nil
}
