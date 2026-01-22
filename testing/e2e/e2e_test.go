package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
	"user-notes-api/auth"
	"user-notes-api/services"

	"github.com/stretchr/testify/assert"
)

var base_url = "http://localhost:8080"

/*
1. Try to get note from other user
2. Get Notes
*/
func TestRegisterLoginCreateNote(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode.")
	}

	waitForServer(t, base_url)

	creds := auth.Credentials{Username: "Alice", Password: "secret_pwd"}
	body, err := json.Marshal(creds)

	if err != nil {
		t.Fatal(err)
	}

	token := callPost(t, base_url, "/register", body)
	assert.True(t, len(token) > 0)

	note := services.Note{Title: "e2e note", Content: "This note is created for the e2e test."}

	body, err = json.Marshal(note)

	if err != nil {
		t.Fatal(err)
	}

	id := callAuthPost(t, base_url, "/notes", token, body)
	assert.True(t, id > 0)

	note_resp, status_code := callGetSingleNote(t, base_url, id, token)
	assert.Equal(t, note, note_resp)
	assert.Equal(t, status_code, http.StatusOK)
}

func TestGetNoteWrongUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode.")
	}

	waitForServer(t, base_url)

	creds := auth.Credentials{Username: "Bob", Password: "secret_pwd"}
	creds_wrong := auth.Credentials{Username: "Clint", Password: "secret_pwd"}
	body, err := json.Marshal(creds)

	if err != nil {
		t.Fatal(err)
	}

	token := callPost(t, base_url, "/register", body)
	assert.True(t, len(token) > 0)

	body, err = json.Marshal(creds_wrong)

	if err != nil {
		t.Fatal(err)
	}

	token_wrong := callPost(t, base_url, "/register", body)
	assert.True(t, len(token_wrong) > 0)

	note := services.Note{Title: "e2e note", Content: "This note is created for the e2e test."}

	body, err = json.Marshal(note)

	if err != nil {
		t.Fatal(err)
	}

	id := callAuthPost(t, base_url, "/notes", token, body)
	assert.True(t, id > 0)

	note_resp, status_code := callGetSingleNote(t, base_url, id, token_wrong)
	assert.Equal(t, services.Note{}, note_resp)
	assert.Equal(t, status_code, http.StatusUnauthorized)
}

func TestGetNotes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode.")
	}

	waitForServer(t, base_url)

	creds := auth.Credentials{Username: "David", Password: "secret_pwd"}

	body, err := json.Marshal(creds)

	if err != nil {
		t.Fatal(err)
	}

	token := callPost(t, base_url, "/register", body)
	assert.True(t, len(token) > 0)

	note1 := services.Note{Title: "e2e note2", Content: "This note is created for the e2e test."}

	body, err = json.Marshal(note1)

	if err != nil {
		t.Fatal(err)
	}

	id1 := callAuthPost(t, base_url, "/notes", token, body)
	assert.True(t, id1 > 0)

	note2 := services.Note{Title: "e2e note2", Content: "This note is created for the e2e test."}

	body, err = json.Marshal(note2)

	if err != nil {
		t.Fatal(err)
	}

	id2 := callAuthPost(t, base_url, "/notes", token, body)
	assert.True(t, id2 > 0)

	notes := callGetNotes(t, base_url, token)
	assert.Equal(t, len(notes.Result), 2)
	assert.Equal(t, id1, notes.Result[0].Id)
	assert.Equal(t, id2, notes.Result[1].Id)
	assert.Equal(t, note1.Title, notes.Result[0].Title)
	assert.Equal(t, note2.Title, notes.Result[1].Title)
}
