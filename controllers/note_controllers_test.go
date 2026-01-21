package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-notes-api/services"
	"user-notes-api/testing/testutils/servicemocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNoteControllerCreateSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	note := services.Note{Title: "title", Content: "content"}
	marshalled, err := json.Marshal(note)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(marshalled))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Set("username", "Alice")

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note_mod_service.On("CreateNote", req_ctx, note, "Alice").Return(1, nil)

	note_controller.Create(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
}

func TestNoteControllerUsernameWrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	note := services.Note{Title: "title", Content: "content"}
	marshalled, err := json.Marshal(note)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(marshalled))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Set("username", 5)

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note_mod_service.On("CreateNote", req_ctx, note, "Alice").Return(1, nil)

	note_controller.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `something went wrong when parsing username from token`)
}

func TestNoteControllerMissingUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	note := services.Note{Title: "title", Content: "content"}
	marshalled, err := json.Marshal(note)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(marshalled))
	c.Request.Header.Set("Content-Type", "application/json")

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note_mod_service.On("CreateNote", req_ctx, note, "Alice").Return(1, nil)

	note_controller.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `failed to parse username from context`)
}

func TestNoteControllerCreatorServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	note := services.Note{Title: "title", Content: "content"}
	marshalled, err := json.Marshal(note)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(marshalled))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Set("username", "Alice")

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note_mod_service.On("CreateNote", req_ctx, note, "Alice").Return(0, &services.ErrorUserNotFound{Username: "Alice"})

	note_controller.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `user Alice not found`)
}

func TestNoteControllerGetNotesSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notes", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	c.Set("user_id", uint(1))

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	var notes services.GetNotesResult
	notes.Result = append(notes.Result, services.NoteListResult{Id: 1, Title: "Title1"})
	notes.Result = append(notes.Result, services.NoteListResult{Id: 2, Title: "Title2"})
	note_read_service.On("GetNotes", req_ctx, uint(1)).Return(notes, nil)

	note_controller.GetNotes(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Id":1`)
	assert.Contains(t, w.Body.String(), `"Title":"Title1"`)
	assert.Contains(t, w.Body.String(), `"Id":2`)
	assert.Contains(t, w.Body.String(), `"Title":"Title2"`)
}

func TestNoteControllerGetSingleNoteSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notes/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set("user_id", uint(1))

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note := services.Note{Title: "Test title", Content: "Test content"}
	note_read_service.On("GetNote", req_ctx, uint(1), uint(1)).Return(note, nil)

	note_controller.GetSingleNote(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Title":"Test title"`)
	assert.Contains(t, w.Body.String(), `"Content":"Test content"`)
}

func TestNoteControllerGetSingleNoteWrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notes/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set("user_id", uint(1))

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note := services.Note{}
	e := services.ErrorWrongOwner{UserId: 1, NoteId: 1}
	note_read_service.On("GetNote", req_ctx, uint(1), uint(1)).Return(note, &e)

	note_controller.GetSingleNote(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "does not own note")
}

func TestNoteControllerGetSingleNoteNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notes/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Set("user_id", uint(1))

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	note := services.Note{}
	e := services.ErrorNoteNotFound{NoteId: 1}
	note_read_service.On("GetNote", req_ctx, uint(1), uint(1)).Return(note, &e)

	note_controller.GetSingleNote(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestNoteControllerGetNotesNotesNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/notes", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	c.Set("user_id", uint(1))

	note_mod_service := new(servicemocks.MockNoteModificationService)
	note_read_service := new(servicemocks.MockNoteReaderService)
	note_controller := NewNoteController(note_mod_service, note_read_service)

	req_ctx := c.Request.Context()
	var notes services.GetNotesResult
	e := services.ErrorNotesNotFound{UserId: 1, Err: errors.New("user not found")}
	note_read_service.On("GetNotes", req_ctx, uint(1)).Return(notes, &e)

	note_controller.GetNotes(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
