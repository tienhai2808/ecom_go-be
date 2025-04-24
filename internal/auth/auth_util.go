package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string) (string, error) {
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

func verifyPassword(storedPassword, inputPassword string) (bool, error) {

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