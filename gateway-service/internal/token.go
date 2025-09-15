package internal

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenExpired = errors.New("TOKEN_EXPIRED")
var ErrTokenInvalid = errors.New("TOKEN_INVALID")
var ErrTokenEmpty = errors.New("TOKEN_EMPTY")

type TokenInfo struct {
	UserID string
	Roles  []string
}

func ExtractTokenInfo(tokenString string) (*TokenInfo, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrTokenEmpty
	}

	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		return nil, errors.New("missing AUTH_JWT_SECRET")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))

	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	userID := toString(claims["user_id"])
	roles := extractRoles(claims["roles"])

	if t, err := claims.GetExpirationTime(); err == nil && t != nil {
		if time.Now().UTC().After(t.Time) {
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
