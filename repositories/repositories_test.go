package repositories

import (
	"context"
	"testing"

	"user-notes-api/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func prepareDatabase(t *testing.T) (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to SQLite db:", err)
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Note{})

	return db
}

func TestUserRepository(t *testing.T) {
	db := prepareDatabase(t)
	ctx := context.Background()

	userRepo := UserRepository{db: db}
	// Create User via User object
	user := models.User{Username: "Alice", Password: "pwd"}
	err := userRepo.CreateUser(ctx, &user)

	assert.NoError(t, err)
	assert.NotEqual(t, user.ID, int64(0))

	// Create User via username and password
	user2, err := userRepo.CreateUserByNameAndPassword(ctx, "Bob", "hashed")
	assert.NoError(t, err)
	assert.NotEqual(t, user2.ID, int64(0))
	assert.Equal(t, "Bob", user2.Username)
	assert.Equal(t, "hashed", user2.Password)

	// Cannot create users with the same name
	user_duplicate := models.User{Username: "Alice", Password: "pwd"}
	err = userRepo.CreateUser(ctx, &user_duplicate)
	assert.Error(t, err)

	user_duplicate, err = userRepo.CreateUserByNameAndPassword(ctx, "Alice", "hashed")
	assert.Error(t, err)

	// Find User by Id
	user_read, err := userRepo.findUserById(ctx, user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", user_read.Username)
	assert.Equal(t, "hashed", user_read.Password)
	assert.Equal(t, user2.ID, user_read.ID)

	// Find User by name
	user_read, err = userRepo.findUserByName(ctx, "Alice")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user_read.Username)
	assert.Equal(t, "pwd", user_read.Password)
	assert.Equal(t, user.ID, user_read.ID)

	// Error when trying to find non existent user
	user_read, err = userRepo.findUserByName(ctx, "Clint")
	assert.Error(t, err)

	id := max(user.ID, user2.ID) + 1
	user_read, err = userRepo.findUserById(ctx, id)
	assert.Error(t, err)

	// Update user
	// Delete user via Id
	// Delete user via User object

}

func TestNoteRepository(t *testing.T) {
	// Find by Id
	// Create via Note object
	// Update
	// Delete via id

}

func TestCascadingDelete(t *testing.T) {

}
