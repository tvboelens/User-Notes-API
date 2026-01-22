package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
	"user-notes-api/services"
)

type JwtToken struct {
	Token string `json:"token"`
}

type IdField struct {
	ID uint `json:"id"`
}

func waitForServer(t *testing.T, baseUrl string) {
	for range 40 {
		resp, err := http.Get(baseUrl + "/health")
		if err == nil && resp.StatusCode == 200 {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("server never became ready")
}

func callPost(t *testing.T, base_url string, path string, body []byte) string {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", base_url+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Response code not OK: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var token JwtToken
	err = json.Unmarshal(resp_body, &token)
	if err != nil {
		t.Fatal(err)
	}

	return token.Token
}

func callAuthPost(t *testing.T, base_url string, path string, jwt_token string, body []byte) uint {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", base_url+path, bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+jwt_token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Response code not OK: " + strconv.Itoa(resp.StatusCode))
	}

	resp_body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var id_field IdField
	err = json.Unmarshal(resp_body, &id_field)
	if err != nil {
		t.Fatal(err)
	}

	return id_field.ID
}

func callGetSingleNote(t *testing.T, base_url string, note_id uint, jwt_token string) (services.Note, int) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", base_url+"/notes/"+strconv.Itoa(int(note_id)), nil)
	req.Header.Add("Authorization", "Bearer "+jwt_token)

	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return services.Note{}, resp.StatusCode
	}
	resp_body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var note services.Note
	err = json.Unmarshal(resp_body, &note)
	if err != nil {
		t.Fatal(err)
	}

	return note, http.StatusOK
}

func callGetNotes(t *testing.T, base_url string, jwt_token string) services.GetNotesResult {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", base_url+"/notes/", nil)
	req.Header.Add("Authorization", "Bearer "+jwt_token)

	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Response code not OK: " + strconv.Itoa(resp.StatusCode))
	}
	resp_body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var note_result services.GetNotesResult
	err = json.Unmarshal(resp_body, &note_result)
	if err != nil {
		t.Fatal(err)
	}

	return note_result
}
