package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
)

func GenerateGuestToken(guestID string, ttl time.Duration, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": guestID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateToken(userID int64, userRole string, ttl time.Duration, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": userRole,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ExtractGuestToken(claims jwt.MapClaims) (string, error) {
	guestID, ok := claims["sub"].(string)
	if !ok {
		return "", customErr.ErrGuestIdNotFound
	}
	
	return guestID, nil
}

func ParseToken(tokenStr, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, customErr.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, customErr.ErrInvalidToken
}

func ExtractToken(claims jwt.MapClaims) (int64, string, error) {
	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", customErr.ErrUserIdNotFound
	}

	userID := int64(userIDFloat)

	userRole, ok := claims["role"].(string)
	if !ok {
		return 0, "", customErr.ErrUserRoleNotFound
	}

	return userID, userRole, nil
}
