package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	mrand "math/rand"
	"strings"

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
		return "", fmt.Errorf("lỗi khi tạo salt băm mật khẩu: %w", err)
	}

	timeCost := uint32(3)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, timeCost, memory, threads, keyLength)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, timeCost, threads, encodedSalt, encodedHash)

	return encoded, nil
}

func VerifyPassword(storedPassword, inputPassword string) (bool, error) {
	parts := strings.Split(storedPassword, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, fmt.Errorf("mật khẩu đã lưu không đúng định dạng")
	}

	var memory, timeCost uint32
	var threads uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &timeCost, &threads)
	if err != nil {
		return false, fmt.Errorf("lỗi khi đọc tham số: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) != 16 {
		return false, fmt.Errorf("lỗi khi giải mã salt: %w", err)
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("lỗi khi giải mã hash: %w", err)
	}

	inputHash := argon2.IDKey([]byte(inputPassword), salt, timeCost, memory, threads, uint32(len(storedHash)))

	return subtle.ConstantTimeCompare(storedHash, inputHash) == 1, nil
}