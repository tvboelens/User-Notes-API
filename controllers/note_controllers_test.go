package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-notes-api/services"
	"user-notes-api/testing/testutils/servicemocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNoteControllerSuccess(t *testing.T) {
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
