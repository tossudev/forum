package password

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"hash"
)

// GenerateSalt generates a salt (random value) of provided length (byte size).
// Add salt to the password for secure hashing.
// Minimum recommended salt length: 16 bytes.
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}
	return salt, nil
}

// HashPassword derives a 32-byte key ("hash") using PBKDF2
// with hashing algorithm SHA512 and 10000 iterations.
func HashPassword(password string, salt []byte) ([]byte, error) {
	// PBKDF2 parameters
	params := getHashParams()

	// Derive key using PBKDF2
	key, err := pbkdf2.Key(params.hashAlgo, password, salt, params.iterations, params.keyLength)
	if err != nil {
		return nil, fmt.Errorf("deriving PBKDF2 key: %w", err)
	}
	return key, nil
}

// VerifyPassword checks that the hash of the input password matches the stored hash
// by hashing the incoming password with the original salt
func VerifyPassword(inputPassword string, storedHash string, storedSalt []byte) (bool, error) {
	// Decode the stored storedHash from base64 string into bytes
	decodedHash, _ := base64.StdEncoding.DecodeString(storedHash)

	// PBKDF2 parameters
	params := getHashParams()

	hashedPassword, err := pbkdf2.Key(params.hashAlgo, inputPassword, storedSalt, params.iterations, params.keyLength)
	if err != nil {
		return false, fmt.Errorf("deriving PBKDF2 key: %w", err)
	}

	// Compare the hashes
	return subtle.ConstantTimeCompare(decodedHash, hashedPassword) == 1, nil
}

type hashParams struct {
	iterations int
	keyLength  int // Number of bytes
	hashAlgo   func() hash.Hash
}

func getHashParams() hashParams {
	params := hashParams{
		iterations: 10000,
		keyLength:  32,
		hashAlgo:   sha512.New,
	}
	return params
}
