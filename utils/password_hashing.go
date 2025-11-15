package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"

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

// hash string encoding according to https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
func EncodeHashString(ph *ParsedHashString) (string, error) {
	hash_string := "$" + ph.Id

	if ph.Version > 0 {
		hash_string += "$v=" + strconv.Itoa(ph.Version)
	}

	if len(ph.Params) > 0 {
		hash_string += "$"
		for key, val := range ph.Params {
			hash_string += key + "=" + strconv.FormatUint(uint64(val), 10) + ","
		}
		// remove trailing comma
		hash_string = strings.TrimSuffix(hash_string, ",")
	}

	if len(ph.Salt) > 0 {
		encoded_salt := "$" + base64.RawStdEncoding.EncodeToString(ph.Salt)
		hash_string += encoded_salt

		// hash may only be present if a salt is present
		if len(ph.Hash) > 0 {
			encoded_hash := "$" + base64.RawStdEncoding.EncodeToString(ph.Hash)
			hash_string += encoded_hash
		}
	}

	return hash_string, nil
}

func ParseHashString(hash_string string) (ParsedHashString, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(hash_string)
	ph := ParsedHashString{Hash: decoded}
	return ph, err
}
