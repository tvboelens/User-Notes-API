package services

import (
	"context"
	"fmt"

	"user-notes-api/models"
	"user-notes-api/repositories"
)

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteReaderService interface {
	getNotes(ctx context.Context, userId uint) ([]Note, error)
	getNote(ctx context.Context, noteId uint, userId uint) (Note, error)
}

type NoteModificationService interface {
	CreateNote(ctx context.Context, note Note, username string) (uint, error)
	UpdateNote(ctx context.Context, id uint, note Note) error
	DeleteNote(ctx context.Context, id uint) error
}

type ErrorUserNotFound struct {
	Username string
	Err      error
}

type ErrorNoteNotFound struct {
	NoteId uint
	Err    error
}

type ErrorNotesNotFound struct {
	UserId uint
	Err    error
}

type ErrorWrongOwner struct {
	NoteId uint
	UserId uint
}

func (e *ErrorUserNotFound) Error() string {
	return fmt.Sprintf("user %q not found: %v", e.Username, e.Err)
}

func (e *ErrorUserNotFound) Unwrap() error {
	return e.Err
}

func (e *ErrorNoteNotFound) Error() string {
	return fmt.Sprintf("note with id %d not found: %v", e.NoteId, e.Err)
}

func (e *ErrorNoteNotFound) Unwrap() error {
	return e.Err
}

func (e *ErrorNotesNotFound) Error() string {
	return fmt.Sprintf("notes from user with id %d not found: %v", e.UserId, e.Err)
}

func (e *ErrorNotesNotFound) Unwrap() error {
	return e.Err
}

func (e *ErrorWrongOwner) Error() string {
	return fmt.Sprintf("user with id %d does not own note with id %d", e.UserId, e.NoteId)
}

type NoteService struct {
	UserRepo    repositories.UserReader
	NoteCreator repositories.NoteCreator
	NoteReader  repositories.NoteReader
}

func NewNoteService(note_reader repositories.NoteReader, note_creator repositories.NoteCreator, user_repo repositories.UserReader) *NoteService {
	note_service := NoteService{NoteReader: note_reader, NoteCreator: note_creator, UserRepo: user_repo}
	return &note_service
}

func (s *NoteService) getNote(ctx context.Context, noteId uint, userId uint) (Note, error) {
	note, err := s.NoteReader.FindNoteById(ctx, noteId)

	if err != nil {
		return Note{}, &ErrorNoteNotFound{NoteId: noteId, Err: err}
	}

	if note.UserID != userId {
		return Note{}, &ErrorWrongOwner{NoteId: noteId, UserId: userId}
	}

	return Note{Title: note.Title, Content: note.Body}, nil
}

func (s *NoteService) getNotes(ctx context.Context, userId uint) ([]Note, error) {
	var note_array []Note
	notes, err := s.NoteReader.FindNotesByUserId(ctx, userId)
	if err != nil {
		return note_array, &ErrorNotesNotFound{UserId: userId, Err: err}
	}

	for _, note := range *notes {
		if note.UserID != userId {
			return note_array, &ErrorWrongOwner{NoteId: note.ID, UserId: userId}
		}
		note_array = append(note_array, Note{Title: note.Title, Content: note.Body})
	}
	return note_array, nil
}

func (s *NoteService) CreateNote(ctx context.Context, note Note, username string) (uint, error) {
	user, err := s.UserRepo.FindUserByName(ctx, username)
	if err != nil {
		return 0, &ErrorUserNotFound{Username: username, Err: err}
	}

	note_model := models.Note{User: *user, UserID: user.ID, Title: note.Title, Body: note.Content}
	err = s.NoteCreator.CreateNote(ctx, &note_model)

	return note_model.ID, err
}
func (s *NoteService) UpdateNote(ctx context.Context, id uint, note Note) error {
	return nil
}
func (s *NoteService) DeleteNote(ctx context.Context, id uint) error {
	return nil
}
