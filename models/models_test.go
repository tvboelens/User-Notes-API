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

	db.AutoMigrate(&User{})

	return db
}

func TestCreateAndRetrieveUser(t *testing.T) {
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

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()

}
