package utils

import (
	"bytes"
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

	// Verification only works with the correct salt
	wrong_salt, err := hasher.GenerateSalt()
	assert.NoError(t, err)
	assert.NotNil(t, wrong_salt)

	isEqual, err = compare(hash, password, wrong_salt, &hasher)
	assert.NoError(t, err)
	assert.False(t, isEqual)

	// Verification fails with the wrong password
	wrong_pw := []byte("wrong_password")
	isEqual, err = compare(hash, wrong_pw, salt, &hasher)
	assert.NoError(t, err)
	assert.False(t, isEqual)

	isEqual, err = compare(hash, wrong_pw, wrong_salt, &hasher)
	assert.NoError(t, err)
	assert.False(t, isEqual)
}

func TestHashStringEncoding(t *testing.T) {
	hash := []byte("super_long_password")
	salt := []byte("random_salt")

	// encoding hash without salt results in error
	ph := ParsedHashString{Id: "argon2id", Hash: hash}
	_, err := EncodeHashString(&ph)
	assert.Error(t, err)

	// encoding without id results in error
	ph = ParsedHashString{Hash: hash, Salt: salt}
	_, err = EncodeHashString(&ph)
	assert.Error(t, err)

	// success with id, hash and salt
	ph = ParsedHashString{Id: "argon2id", Hash: hash, Salt: salt}
	str, err := EncodeHashString(&ph)

	assert.NoError(t, err)

	decoded_ph, err := ParseHashString(str)

	assert.NoError(t, err)
	assert.True(t, bytes.Equal(hash, decoded_ph.Hash))
	assert.True(t, bytes.Equal(salt, decoded_ph.Salt))
	assert.Equal(t, "argon2id", decoded_ph.Id)

	params := make(map[string]uint32)
	params["Time"] = 1
	params["SaltLen"] = 32
	params["Memory"] = 64 * 1024

	// success with all args
	ph = ParsedHashString{Hash: hash, Salt: salt, Id: "argon2id", Version: 19, Params: params}

	str, err = EncodeHashString(&ph)
	assert.NoError(t, err)

	decoded_ph, err = ParseHashString(str)

	assert.NoError(t, err)
	assert.True(t, bytes.Equal(hash, decoded_ph.Hash))
	assert.True(t, bytes.Equal(salt, decoded_ph.Salt))
	assert.Equal(t, "argon2id", decoded_ph.Id)
	assert.Equal(t, 19, decoded_ph.Version)
	assert.Equal(t, params, decoded_ph.Params)

	// success without params
	ph = ParsedHashString{Hash: hash, Salt: salt, Id: "argon2id", Version: 19}

	str, err = EncodeHashString(&ph)
	assert.NoError(t, err)

	decoded_ph, err = ParseHashString(str)

	assert.NoError(t, err)
	assert.True(t, bytes.Equal(hash, decoded_ph.Hash))
	assert.True(t, bytes.Equal(salt, decoded_ph.Salt))
	assert.Equal(t, "argon2id", decoded_ph.Id)
	assert.Equal(t, 19, decoded_ph.Version)

	// success without version but with params
	ph = ParsedHashString{Hash: hash, Salt: salt, Id: "argon2id", Params: params}

	str, err = EncodeHashString(&ph)
	assert.NoError(t, err)

	decoded_ph, err = ParseHashString(str)

	assert.NoError(t, err)
	assert.True(t, bytes.Equal(hash, decoded_ph.Hash))
	assert.True(t, bytes.Equal(salt, decoded_ph.Salt))
	assert.Equal(t, "argon2id", decoded_ph.Id)
	assert.Equal(t, params, decoded_ph.Params)
}
