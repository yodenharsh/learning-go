package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type ResetPasswordCode struct {
	Value     string
	ExpiresAt time.Time
}

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

func CreatePasswordResetToken() (ResetPasswordCode, error) {
	parsedDuration, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_TOKEN_EXP_DURATION_IN_MINUTES"))

	if err != nil {
		return ResetPasswordCode{}, ErrorHandler(err, "Failed when reading and parsing env variable")
	}
	mins := time.Duration(parsedDuration)
	expiry := time.Now().Add(mins * time.Minute)

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return ResetPasswordCode{}, ErrorHandler(err, "Error during creating random bytes")
	}

	hashedToken := sha256.Sum256(tokenBytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	return ResetPasswordCode{
		Value:     hashedTokenString,
		ExpiresAt: expiry,
	}, err
}
