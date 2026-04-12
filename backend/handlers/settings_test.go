package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSettingsPUT(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	test_settings := struct {
		email    string
		password string
	}{
		email:    "test@example.com",
		password: "abc123",
	}

	settingsJSON, _ := json.Marshal(test_settings)
	req := httptest.NewRequest("PUT", "/api/settings", bytes.NewBuffer(settingsJSON))
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	UpdateSettings(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for PUT /api/settings: Expected 200, got %d", w.Result().StatusCode)
	}
}

func TestSettingsGET(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	req := httptest.NewRequest("GET", "/api/settings", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	GetSettings(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for GET /api/settings: Expected 200, got %d", w.Result().StatusCode)
	}

	test_settings := struct {
		email    string
		password string
	}{}

	if err := json.NewDecoder(w.Body).Decode(&test_settings); err != nil {
		t.Fatalf("Failed to decode settings response: %v", err)
	}
}
