package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gofiber/fiber/v2/log"
)

// GenerateSixDigitPIN generates a random 6-digit PIN
func GenerateSixDigitPIN() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Errorf("GenerateSixDigitPIN: Failed to generate six digit PIN: %v", err)
		return ""
	}
	return fmt.Sprintf("%06d", n.Int64())
}
