package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "43041")
	os.Setenv("DB_NAME", "UserNotesAPI_DB")
	os.Setenv("JWT_SECRET", "JWTsecret")

	cfg := LoadConfig()

	assert.Equal(t, "root", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "UserNotesAPI_DB", cfg.DBName)
	assert.Equal(t, "43041", cfg.DBPort)
	assert.Equal(t, "JWTsecret", cfg.JWTSecret)

}
