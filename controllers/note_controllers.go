package controllers

import (
	"errors"
	"net/http"
	"strconv"

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

func (n *NoteController) GetNotes(c *gin.Context) {
	request_ctx := c.Request.Context()
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse user id from context"})
		return
	}

	user_id, ok := uid.(uint)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert user id to uint"})
		return
	}

	result, err := n.ReaderService.GetNotes(request_ctx, user_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notes not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *NoteController) GetSingleNote(c *gin.Context) {
	request_ctx := c.Request.Context()
	note_id_str := c.Param("id")
	note_id, err := strconv.Atoi(note_id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed id"})
		return
	}

	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user id from context"})
		return
	}

	user_id, ok := uid.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed user id"})
		return
	}

	note, err := n.ReaderService.GetNote(request_ctx, uint(note_id), user_id)
	if err != nil {
		var e *services.ErrorWrongOwner
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, note)
}
