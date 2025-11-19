package repositories

import (
	"context"
	"errors"
	"fmt"
	"user-notes-api/models"

	"gorm.io/gorm"
)

type NoteReader interface {
	FindNoteById(ctx context.Context, id uint) (*models.Note, error)
	FindNotesByUserId(ctx context.Context, userId uint) (*[]models.Note, error)
}

type NoteCreator interface {
	CreateNote(ctx context.Context, note *models.Note) error
}

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	tx := r.db.WithContext(ctx).Omit("User").Create(note)

	if tx.Error == nil && tx.RowsAffected != 1 {
		return errors.New("number of affected rows not equal to 1")
	}

	return tx.Error
}

func (r *NoteRepository) FindNoteById(ctx context.Context, id uint) (*models.Note, error) {
	note, err := gorm.G[models.Note](r.db).Where("id = ?", id).First(ctx)
	return &note, err
}

func (r *NoteRepository) FindNotesByUserId(ctx context.Context, userId uint) (*[]models.Note, error) {
	notes, err := gorm.G[models.Note](r.db).Where("user_id = ?", userId).Find(ctx)
	return &notes, err
}

func (r *NoteRepository) DeleteNote(ctx context.Context, note *models.Note) error {
	count, err := gorm.G[models.Note](r.db).Where("id = ?", note.ID).Delete(ctx)
	if err == nil && count != 1 {
		msg := fmt.Sprintf("unexpected count for deleting note. expected 1, received %d", count)
		return errors.New(msg)
	}
	return err

}

func (r *NoteRepository) DeleteNoteById(ctx context.Context, id uint) error {
	count, err := gorm.G[models.Note](r.db).Where("id = ?", id).Delete(ctx)
	if err == nil && count != 1 {
		msg := fmt.Sprintf("unexpected count for deleting note. expected 1, received %d", count)
		return errors.New(msg)
	}
	return err

}

func (r *NoteRepository) DeleteNotesOfUser(ctx context.Context, user *models.User) error {
	no_of_notes := len(user.Notes)
	count, err := gorm.G[models.Note](r.db).Where("user_id = ?", user.ID).Delete(ctx)
	if err == nil && count != no_of_notes {
		msg := fmt.Sprintf("unexpected count for deleting notes. expected %d, received %d", no_of_notes, count)
		return errors.New(msg)
	}
	return err

}

func (r *NoteRepository) DeleteNotesOfUserByUserID(ctx context.Context, id uint) (int, error) {
	count, err := gorm.G[models.Note](r.db).Where("user_id = ?", id).Delete(ctx)
	return count, err

}
