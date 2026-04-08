package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bananawafflecookies.com/m/v2/db"
)

func TestJobPOST(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	test_job := db.Job{
		UserID:       test_user.Id,
		CompanyName:  "Aperture Labs",
		Title:        "Tester",
		Status:       "applied",
		DeadlineDate: time.Now(),
		Description:  "Portal Gun tester",
	}

	jobJSON, _ := json.Marshal(test_job)
	req := httptest.NewRequest("POST", "/api/jobs", bytes.NewBuffer(jobJSON))
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	CreateJob(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for POST /api/jobs: Expected 200, got %d", w.Result().StatusCode)
	}
}

func TestJobsGET(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	req := httptest.NewRequest("GET", "/api/jobs", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	GetJobs(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for GET /api/jobs: Expected 200, got %d", w.Result().StatusCode)
	}

	var jobs []db.Job
	if err := json.NewDecoder(w.Body).Decode(&jobs); err != nil {
		t.Fatalf("Failed to decode jobs response: %v", err)
	}
	if len(jobs) == 0 {
		t.Fatal("Expected at least one job in response")
	}
}
