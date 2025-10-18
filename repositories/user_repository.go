package repositories

import (
	"context"
	"errors"
	"user-notes-api/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	result := gorm.WithResult()
	err := gorm.G[models.User](r.db, result).Create(ctx, user)

	if err == nil && result.RowsAffected != 1 {
		return errors.New("number of affected rows not equal to 1")
	}

	return err
}

func (r *UserRepository) CreateUserByNameAndPassword(ctx context.Context, username string, password string) (models.User, error) {
	user := models.User{Username: username, Password: password}
	err := r.CreateUser(ctx, &user)
	return user, err
}
