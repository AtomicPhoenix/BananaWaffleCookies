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
		t.Fatalf("GET /api/profile: Expected 200, got %d", w.Result().StatusCode)
	}

	var profile db.Profile
	if err := json.NewDecoder(w.Body).Decode(&profile); err != nil {
		t.Fatalf("Failed to decode profile response: %v", err)
	}

	if profile.UserID != test_user.Id {
		t.Errorf("Expected UserID %d, got %d", test_user.Id, profile.UserID)
	}
}

func TestProfilePUT(t *testing.T) {
	test_user := createTestUser(t)
	cookie := getAuthCookie(t, test_user.Id, "test@example.com")
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	profile := db.Profile{
		FirstName:          "Test",
		LastName:           "User",
		Phone:              "1234567890",
		Location:           "New Jersey",
		City:               "Newark",
		State:              "NJ",
		Country:            "USA",
		Headline:           "Software Engineer",
		LinkedinURL:        "https://linkedin.com/in/test",
		PortfolioURL:       "https://test.dev",
		Summary:            "Building things.",
		PreferredCity:      "New York",
		PreferredState:     "NY",
		PreferredRole:      "Backend Engineer",
		PreferredSalaryMin: 120000,
		PreferredSalaryMax: 180000,
		WorkMode:           "hybrid",
	}
	profile.SetCompletionPercent()

	profileJSON, _ := json.Marshal(profile)

	req := httptest.NewRequest("PUT", "/api/profile", bytes.NewBuffer(profileJSON))
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	UpdateProfile(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("PUT /api/profile: Expected 200, got %d", w.Result().StatusCode)
	}

	// Verify update
	getReq := httptest.NewRequest("GET", "/api/profile", nil)
	getReq.AddCookie(cookie)

	getW := httptest.NewRecorder()
	GetProfile(getW, getReq)

	if getW.Result().StatusCode != http.StatusOK {
		t.Fatalf("GET after PUT failed: got %d", getW.Result().StatusCode)
	}

	var retrieved_profile db.Profile
	if err := json.NewDecoder(getW.Body).Decode(&retrieved_profile); err != nil {
		t.Fatalf("Failed to decode GET response: %v", err)
	}

	if retrieved_profile.UserID != test_user.Id {
		t.Errorf("UserID: expected %d, got %d", test_user.Id, retrieved_profile.UserID)
	}
	if retrieved_profile.FirstName != profile.FirstName {
		t.Errorf("FirstName: expected %s, got %s", profile.FirstName, retrieved_profile.FirstName)
	}
	if retrieved_profile.LastName != profile.LastName {
		t.Errorf("LastName: expected %s, got %s", profile.LastName, retrieved_profile.LastName)
	}
	if retrieved_profile.Phone != profile.Phone {
		t.Errorf("Phone: expected %s, got %s", profile.Phone, retrieved_profile.Phone)
	}
	if retrieved_profile.Location != profile.Location {
		t.Errorf("Location: expected %s, got %s", profile.Location, retrieved_profile.Location)
	}
	if retrieved_profile.City != profile.City {
		t.Errorf("City: expected %s, got %s", profile.City, retrieved_profile.City)
	}
	if retrieved_profile.State != profile.State {
		t.Errorf("State: expected %s, got %s", profile.State, retrieved_profile.State)
	}
	if retrieved_profile.Country != profile.Country {
		t.Errorf("Country: expected %s, got %s", profile.Country, retrieved_profile.Country)
	}

	if retrieved_profile.Headline != profile.Headline {
		t.Errorf("Headline: expected %s, got %s", profile.Headline, retrieved_profile.Headline)
	}
	if retrieved_profile.LinkedinURL != profile.LinkedinURL {
		t.Errorf("LinkedinURL: expected %s, got %s", profile.LinkedinURL, retrieved_profile.LinkedinURL)
	}
	if retrieved_profile.PortfolioURL != profile.PortfolioURL {
		t.Errorf("PortfolioURL: expected %s, got %s", profile.PortfolioURL, retrieved_profile.PortfolioURL)
	}
	if retrieved_profile.Summary != profile.Summary {
		t.Errorf("Summary: expected %s, got %s", profile.Summary, retrieved_profile.Summary)
	}

	if retrieved_profile.PreferredCity != profile.PreferredCity {
		t.Errorf("PreferredCity: expected %s, got %s", profile.PreferredCity, retrieved_profile.PreferredCity)
	}
	if retrieved_profile.PreferredState != profile.PreferredState {
		t.Errorf("PreferredState: expected %s, got %s", profile.PreferredState, retrieved_profile.PreferredState)
	}
	if retrieved_profile.PreferredRole != profile.PreferredRole {
		t.Errorf("PreferredRole: expected %s, got %s", profile.PreferredRole, retrieved_profile.PreferredRole)
	}
	if retrieved_profile.PreferredSalaryMin != profile.PreferredSalaryMin {
		t.Errorf("PreferredSalaryMin: expected %d, got %d", profile.PreferredSalaryMin, retrieved_profile.PreferredSalaryMin)
	}
	if retrieved_profile.PreferredSalaryMax != profile.PreferredSalaryMax {
		t.Errorf("PreferredSalaryMax: expected %d, got %d", profile.PreferredSalaryMax, retrieved_profile.PreferredSalaryMax)
	}
	if retrieved_profile.WorkMode != profile.WorkMode {
		t.Errorf("WorkMode: expected %s, got %s", profile.WorkMode, retrieved_profile.WorkMode)
	}

	if retrieved_profile.CompletionPercent != profile.CompletionPercent {
		t.Errorf("CompletionPercent: expected %d, got %d", profile.CompletionPercent, retrieved_profile.CompletionPercent)
	}
}
