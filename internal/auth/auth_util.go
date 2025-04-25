package auth

import (
	"crypto/subtle"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/argon2"
)

func GenerateOtp(length int) string {
	min := int(math.Pow10(length))
	max := 9 * min  
	otp := min + mrand.Intn(max)
	return fmt.Sprintf("%d", otp)
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", nil
	}

	timeCost := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, timeCost, memory, threads, keyLength)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := encodedSalt + "." + encodedHash

	return encoded, nil
}

func VerifyPassword(storedPassword, inputPassword string) (bool, error) {

	parts := strings.Split(storedPassword, ".")
	if len(parts) != 2 {
		return false, errors.New("invalid stored password format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	timeCost := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)

	inputHash := argon2.IDKey([]byte(inputPassword), salt, timeCost, memory, threads, keyLength)

	return subtle.ConstantTimeCompare(storedHash, inputHash) == 1, nil
}

func GenerateToken(userID string, userRole string, ttl time.Duration, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    userRole,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}