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
	user_read, err := userRepo.FindUserById(ctx, user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", user_read.Username)
	assert.Equal(t, "hashed", user_read.Password)
	assert.Equal(t, user2.ID, user_read.ID)

	// Find User by name
	user_read, err = userRepo.FindUserByName(ctx, "Alice")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user_read.Username)
	assert.Equal(t, "pwd", user_read.Password)
	assert.Equal(t, user.ID, user_read.ID)

	// Error when trying to find non existent user
	user_read, err = userRepo.FindUserByName(ctx, "Clint")
	assert.Error(t, err)

	id := max(user.ID, user2.ID) + 1
	user_read, err = userRepo.FindUserById(ctx, id)
	assert.Error(t, err)

	// Update user
	// Delete user via Id
	id = user.ID
	count, err := userRepo.DeleteUserById(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	_, err = userRepo.FindUserById(ctx, id)
	assert.Error(t, err)
	// Delete user via User object
	id = user2.ID
	err = userRepo.DeleteUser(ctx, &user2)
	assert.NoError(t, err)

	_, err = userRepo.FindUserById(ctx, id)
	assert.Error(t, err)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()

}

func TestNoteRepository(t *testing.T) {
	db := prepareDatabase(t)
	ctx := context.Background()

	userRepo := UserRepository{db: db}
	noteRepo := NoteRepository{db: db}

	// Create User
	user := models.User{Username: "Alice", Password: "pwd"}
	user_not_in_db := models.User{Username: "Bob", Password: "pwd"}
	err := userRepo.CreateUser(ctx, &user)

	assert.NoError(t, err)
	assert.NotEqual(t, user.ID, int64(0))

	// Create note
	note1 := models.Note{Title: "Title1", Body: "body1", UserID: user.ID, User: user}
	note2 := models.Note{Title: "Title2", Body: "body2"}
	err = noteRepo.CreateNote(ctx, &note1)
	assert.NoError(t, err)

	// Assert that note cannot be created if the owner is not in DB
	_, err = userRepo.FindUserById(ctx, user_not_in_db.ID)
	assert.Error(t, err)
	_, err = userRepo.FindUserByName(ctx, "Bob")
	assert.Error(t, err)

	err = noteRepo.CreateNote(ctx, &note2)
	assert.Error(t, err)
	note2.UserID = user_not_in_db.ID
	note2.User = user_not_in_db
	err = noteRepo.CreateNote(ctx, &note2)
	assert.Error(t, err)

	// Find by Id
	note_read, err := noteRepo.FindNoteById(ctx, note1.ID)
	assert.NoError(t, err)
	assert.Equal(t, note1.ID, note_read.ID)
	assert.Equal(t, "Title1", note_read.Title)
	assert.Equal(t, "body1", note_read.Body)
	assert.Equal(t, note1.UserID, note_read.UserID)
	// Find all notes by one user
	// Find by list of Ids?

	// Update
	// Delete via id
	id := note1.ID
	err = noteRepo.DeleteNoteById(ctx, id)
	assert.NoError(t, err)

	_, err = noteRepo.FindNoteById(ctx, id)
	assert.Error(t, err)

	// create second note
	note2.UserID = user.ID
	note2.User = user
	err = noteRepo.CreateNote(ctx, &note2)
	assert.NoError(t, err)

	// Delete via note object
	id = note2.ID
	err = noteRepo.DeleteNote(ctx, &note2)
	assert.NoError(t, err)

	_, err = noteRepo.FindNoteById(ctx, id)
	assert.Error(t, err)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()

}

func TestCascadingDelete(t *testing.T) {
	db := prepareDatabase(t)
	ctx := context.Background()

	userRepo := UserRepository{db: db}
	noteRepo := NoteRepository{db: db}

	// Create User
	user := models.User{Username: "Alice", Password: "pwd"}
	err := userRepo.CreateUser(ctx, &user)

	assert.NoError(t, err)
	assert.NotEqual(t, user.ID, int64(0))

	// Create notes
	note1 := models.Note{Title: "Title1", Body: "body1", UserID: user.ID, User: user}
	note2 := models.Note{Title: "Title2", Body: "body2", UserID: user.ID, User: user}
	err = noteRepo.CreateNote(ctx, &note1)
	assert.NoError(t, err)
	err = noteRepo.CreateNote(ctx, &note2)
	assert.NoError(t, err)

	id1 := note1.ID
	id2 := note2.ID
	user_id := user.ID

	// Delete user by id
	count, err := userRepo.DeleteUserById(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	_, err = userRepo.FindUserById(ctx, user_id)
	assert.Error(t, err)

	// Assert that notes cannot be found after delete
	_, err = noteRepo.FindNoteById(ctx, id1)
	assert.Error(t, err)

	_, err = noteRepo.FindNoteById(ctx, id2)
	assert.Error(t, err)

	// Recreate User
	user = models.User{Username: "Alice", Password: "pwd"}
	err = userRepo.CreateUser(ctx, &user)

	assert.NoError(t, err)
	assert.NotEqual(t, user.ID, int64(0))

	// Create notes
	note1 = models.Note{Title: "Title1", Body: "body1", UserID: user.ID, User: user}
	note2 = models.Note{Title: "Title2", Body: "body2", UserID: user.ID, User: user}
	err = noteRepo.CreateNote(ctx, &note1)
	assert.NoError(t, err)
	err = noteRepo.CreateNote(ctx, &note2)
	assert.NoError(t, err)

	id1 = note1.ID
	id2 = note2.ID

	// Delete user by id
	err = userRepo.DeleteUser(ctx, &user)
	assert.NoError(t, err)

	// Assert that notes cannot be found after delete
	_, err = noteRepo.FindNoteById(ctx, id1)
	assert.Error(t, err)

	_, err = noteRepo.FindNoteById(ctx, id2)
	assert.Error(t, err)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
}
