package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenExpired = errors.New("TOKEN_EXPIRED")
var ErrTokenInvalid = errors.New("TOKEN_INVALID")
var ErrTokenEmpty = errors.New("TOKEN_EMPTY")

type TokenInfo struct {
	UserID string
	Roles  []string
}

func ExtractTokenInfo(secret string, tokenString string) (*TokenInfo, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrTokenEmpty
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))

	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(
		tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		log.Errorf("ExtractTokenInfo: Failed to parse token: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Errorf("ExtractTokenInfo: Invalid token: %v", token)
		return nil, ErrTokenInvalid
	}

	userID := toString(claims["user_id"])
	roles := extractRoles(claims["roles"])

	if date, err := claims.GetExpirationTime(); err == nil && date != nil {
		if time.Now().UTC().After(date.Time) {
			log.Warnf("ExtractTokenInfo: Token expired: %v", token)
			return nil, ErrTokenExpired
		}
	}

	return &TokenInfo{UserID: userID, Roles: roles}, nil
}

func extractRoles(v any) []string {
	if v == nil {
		return nil
	}
	switch vv := v.(type) {
	case []string:
		return vv
	case []any:
		out := make([]string, 0, len(vv))
		for _, it := range vv {
			if s := strings.TrimSpace(toString(it)); s != "" {
				out = append(out, s)
			}
		}
		return out
	case string:
		parts := strings.Split(vv, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return out
	default:
		return nil
	}
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strings.TrimSuffix(strings.TrimSuffix(strings.TrimRight(strings.TrimRight(fmtFloat(t), "0"), "."), ".0"), ".00")
	case int64:
		return fmtInt(t)
	case int:
		return fmtInt(int64(t))
	default:
		return ""
	}
}

func fmtFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func fmtInt(i int64) string {
	return strconv.FormatInt(i, 10)
}
