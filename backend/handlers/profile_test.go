package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"bananawafflecookies.com/m/v2/db"
)

func TestProfileGET(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	req := httptest.NewRequest("GET", "/api/profile", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	GetProfile(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for GET /api/profile: Expected 200, got %d", w.Result().StatusCode)
	}

	var profile db.Profile
	if err := json.NewDecoder(w.Body).Decode(&profile); err != nil {
		t.Fatalf("Failed to decode profile response: %v", err)
	}
}

func TestProfilePUT(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	var profile db.Profile
	profile.FirstName = "Test"
	profile.LastName = "User"
	profileJSON, _ := json.Marshal(profile)
	req := httptest.NewRequest("PUT", "/api/profile", bytes.NewBuffer(profileJSON))
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	UpdateProfile(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for PUT /api/profile: Expected 200, got %d", w.Result().StatusCode)
	}
}
