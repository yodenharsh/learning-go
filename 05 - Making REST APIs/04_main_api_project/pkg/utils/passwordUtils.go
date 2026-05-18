package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashAndEncodePassword(originalPassword string) (string, error) {
	if originalPassword == "" {
		return "", ErrorHandler(nil, "Password cannot be empty")
	}

	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(originalPassword), salt, 1, 1024*64, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)

	return encodedHash, nil
}

func CompareHashedEncodedPassword(encodedHash, passwordToCheckAgainst string) (bool, error) {
	if encodedHash == "" || passwordToCheckAgainst == "" {
		return false, ErrorHandler(nil, "Encoded hash and original password cannot be empty")
	}

	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return false, ErrorHandler(nil, "Invalid encoded hash format")
	}

	saltBase64 := parts[0]
	hashBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false, ErrorHandler(err, "Error decoding salt from base64")
	}

	hash, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		return false, ErrorHandler(err, "Error decoding hash from base64")
	}
	computedHashOfPasswordToCheck := argon2.IDKey([]byte(passwordToCheckAgainst), salt, 1, 1024*64, 4, 32)

	return subtle.ConstantTimeCompare(hash, computedHashOfPasswordToCheck) == 1, nil
}
