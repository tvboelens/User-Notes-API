package utils

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateHash(password []byte, hasher PasswordHasher) (hash []byte, salt []byte, err error) {
	salt, err = hasher.GenerateSalt()

	if err != nil {
		return nil, nil, err
	}

	hash, err = hasher.GenerateHash(password, salt)
	return hash, salt, err
}

func compare(hash, password, salt []byte, comparer PasswordComparer) (bool, error) {
	return comparer.Compare(hash, salt, password)
}

func TestPasswordHashing(t *testing.T) {
	password := []byte("password")
	threads := uint8(runtime.GOMAXPROCS(0))

	hasher := Argon2IdHasher{Time: 1, SaltLen: 32, Memory: 64 * 1024, Threads: threads, KeyLen: 256}
	hash, salt, err := generateHash(password, &hasher)
	assert.NoError(t, err)
	assert.NotNil(t, hash)
	isEqual, err := compare(hash, password, salt, &hasher)
	assert.True(t, isEqual)
	assert.NoError(t, err)
}
