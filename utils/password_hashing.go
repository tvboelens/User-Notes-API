package utils

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
	GenerateHash(password, salt []byte) ([]byte, error)
	GenerateSalt() ([]byte, error)
}

type PasswordComparer interface {
	Compare(hash, salt, password []byte) (bool, error)
}

type ParsedHashString struct {
	Hash    []byte
	Salt    []byte
	Id      string
	Version int
	Params  map[string]uint32
}

type Argon2IdHasher struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func (h *Argon2IdHasher) GenerateHash(password, salt []byte) ([]byte, error) {
	var err error
	if len(salt) == 0 {
		salt, err = h.GenerateSalt()
	}

	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(password, salt, h.Time, h.Memory, h.Threads, h.KeyLen)
	return hash, nil
}

func (h *Argon2IdHasher) GenerateSalt() ([]byte, error) {
	salt := make([]byte, h.SaltLen)
	_, err := rand.Read(salt)

	if err != nil {

		return nil, err
	}

	return salt, nil
}

func (h *Argon2IdHasher) Compare(hash, salt, password []byte) (bool, error) {
	hashed_pw, err := h.GenerateHash(password, salt)
	if err != nil {
		return false, err
	}

	return (bytes.Equal(hashed_pw, hash)), nil
}

func EncodeHashString(ph *ParsedHashString) (string, error) {
	return "", nil
}

func ParseHashString(hash_string string) (ParsedHashString, error) {
	var ph ParsedHashString
	return ph, nil
}
