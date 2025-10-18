package repositories

import (
	"context"
	"errors"
	"fmt"
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

func (r *UserRepository) findUserById(ctx context.Context, id uint) (models.User, error) {
	user, err := gorm.G[models.User](r.db).Where("id = ?", id).First(ctx)
	return user, err
}

func (r *UserRepository) findUserByName(ctx context.Context, username string) (models.User, error) {
	user, err := gorm.G[models.User](r.db).Where("username = ?", username).First(ctx)
	return user, err
}

func (r *UserRepository) deleteUser(ctx context.Context, user *models.User) error {
	noteRepo := NoteRepository{db: r.db}
	err := noteRepo.deleteNotesOfUser(ctx, user)
	if err != nil {
		return err
	}

	count, err := gorm.G[models.User](r.db).Where("id = ?", user.ID).Delete(ctx)

	if count != 1 {
		msg := fmt.Sprintf("unexpected count for deleting user. expected 1, received %d", count)
		return errors.New(msg)
	}
	return err
}

func (r *UserRepository) deleteUserById(ctx context.Context, id uint) (int, error) {
	noteRepo := NoteRepository{db: r.db}
	count, err := noteRepo.deleteNotesOfUserByUserID(ctx, id)
	if err != nil {
		return count, err
	}

	ccount, err := gorm.G[models.User](r.db).Where("id = ?", id).Delete(ctx)

	if ccount != 1 {
		msg := fmt.Sprintf("unexpected count for deleting user. expected 1, received %d", ccount)
		return count, errors.New(msg)
	}
	return count, err
}
