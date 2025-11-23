package controllers

import (
	"net/http"

	"user-notes-api/services"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	ModificationService services.NoteModificationService
	ReaderService       services.NoteReaderService
}

func NewNoteController(modification_service services.NoteModificationService, reader_service services.NoteReaderService) *NoteController {
	controller := NoteController{ModificationService: modification_service, ReaderService: reader_service}
	return &controller
}

func (n *NoteController) Create(c *gin.Context) {
	var note services.Note
	err := c.Bind(&note)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input" + err.Error()})
		return
	}

	request_ctx := c.Request.Context()
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse username from context"})
		return
	}

	uname, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong when parsing username from token"})
		return
	}

	id, err := n.ModificationService.CreateNote(request_ctx, note, uname)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
