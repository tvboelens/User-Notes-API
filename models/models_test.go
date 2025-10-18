package models

import (
	"context"
	"testing"

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

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Note{})

	return db
}

func TestCRUDUser(t *testing.T) {
	//setup
	db := prepareDatabase(t)

	ctx := context.Background()
	result := gorm.WithResult()

	//create first user
	user := User{Username: "testName", Password: "pwd"}
	err := gorm.G[User](db, result).Create(ctx, &user)

	assert.NoError(t, err)
	assert.Equal(t, result.RowsAffected, int64(1))
	assert.NotEqual(t, user.ID, int64(0))

	// assert that no two users with same name can be created
	user2 := User{Username: "testName", Password: "hashed"}
	err = gorm.G[User](db, result).Create(ctx, &user2)

	assert.Error(t, err)
	assert.Equal(t, result.RowsAffected, int64(0))

	// create second user
	user2 = User{Username: "testName2", Password: "hashed"}
	err = gorm.G[User](db, result).Create(ctx, &user2)

	assert.NoError(t, err)
	assert.Equal(t, result.RowsAffected, int64(1))
	assert.NotEqual(t, user.ID, int64(0))

	var usr, usr2 User
	query_result := db.First(&usr, "username = ?", user.Username)

	assert.NoError(t, query_result.Error)
	assert.Equal(t, user.ID, usr.ID)

	query_result = db.First(&usr2, "username = ?", user2.Username)

	assert.NoError(t, query_result.Error)
	assert.Equal(t, user2.ID, usr2.ID)

	// Update User
	count, err := gorm.G[User](db).Where("username = ?", user.Username).Update(ctx, "password", "hashed")
	assert.Equal(t, count, 1)
	assert.NoError(t, err)
	query_result = db.First(&usr, "username = ?", user.Username)
	assert.NoError(t, query_result.Error)
	assert.Equal(t, user.ID, usr.ID)
	assert.Equal(t, usr.Password, "hashed")

	// Delete User
	count, err = gorm.G[User](db).Where("username = ?", user2.Username).Delete(ctx)
	assert.Equal(t, count, 1)
	assert.NoError(t, err)
	query_result = db.First(&usr, "username = ?", user2.Username)
	assert.Error(t, query_result.Error)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()

}

func TestCRUDNotes(t *testing.T) {
	//setup
	db := prepareDatabase(t)
	ctx := context.Background()
	result := gorm.WithResult()

	// Note can only be created when user exists
	note := Note{Title: "Title", Body: "Body", UserID: 1}
	err := gorm.G[Note](db, result).Create(ctx, &note)
	assert.Error(t, err)

	// Create and read
	user := User{Username: "testName", Password: "pwd"}
	err = gorm.G[User](db, result).Create(ctx, &user)

	assert.NoError(t, err)
	assert.Equal(t, result.RowsAffected, int64(1))
	assert.NotEqual(t, user.ID, int64(0))

	note = Note{Title: "Title", Body: "Body", UserID: user.ID}
	err = gorm.G[Note](db, result).Create(ctx, &note)

	assert.NoError(t, err)
	assert.Equal(t, result.RowsAffected, int64(1))
	assert.NotEqual(t, note.ID, int64(0))

	// Read from DB
	var note_read Note
	query_result := db.Preload("User").First(&note_read, "id = ?", note.ID)
	assert.NoError(t, query_result.Error)
	assert.Equal(t, note_read.UserID, user.ID)
	assert.Equal(t, note_read.User.ID, user.ID)
	assert.Equal(t, note_read.User.Username, "testName")
	assert.Equal(t, note_read.ID, note.ID)
	assert.Equal(t, note_read.Title, note.Title)

	// Update
	note2 := Note{Title: "Title2", Body: "Body2", UserID: user.ID}
	err = gorm.G[Note](db, result).Create(ctx, &note2)

	assert.NoError(t, err)
	assert.Equal(t, result.RowsAffected, int64(1))
	assert.NotEqual(t, note2.ID, int64(0))

	count, err := gorm.G[Note](db).Where("title = ?", note.Title).Update(ctx, "title", "updatedTitle")
	assert.Equal(t, count, 1)
	assert.NoError(t, err)
	query_result = db.First(&note_read, "id = ?", note.ID)
	assert.NoError(t, query_result.Error)
	assert.Equal(t, note.ID, note_read.ID)
	assert.Equal(t, note_read.Title, "updatedTitle")

	// Delete
	count, err = gorm.G[Note](db).Where("id = ?", note2.ID).Delete(ctx)
	assert.Equal(t, count, 1)
	assert.NoError(t, err)
	query_result = db.First(&note_read, "id = ?", note2.ID)
	assert.Error(t, query_result.Error)

	// Get all notes from user
	var usr User
	note2 = Note{Title: "Title2", Body: "Body2", UserID: user.ID}
	err = gorm.G[Note](db, result).Create(ctx, &note2)

	assert.NoError(t, err)
	assert.Equal(t, result.RowsAffected, int64(1))
	assert.NotEqual(t, note2.ID, int64(0))

	query_result = db.Preload("Notes").First(&usr, "ID = ?", user.ID)
	assert.NoError(t, query_result.Error)
	assert.Equal(t, usr.ID, user.ID)
	assert.Equal(t, len(usr.Notes), 2)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
}
