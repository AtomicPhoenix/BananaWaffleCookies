package db

import "testing"

func TestUpdateProfile(t *testing.T) {
	// Creating a test user also creates a profile automatically
	test_user := createTestUser(t)
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	profile := Profile{
		UserID:       test_user.Id,
		FirstName:    "Cave",
		LastName:     "Johnson",
		Phone:        "5132737724",
		City:         "Unknown",
		State:        "Unknown",
		Country:      "USA",
		LinkedinURL:  "https://half-life.fandom.com/wiki/Cave_Johnson",
		PortfolioURL: "https://half-life.fandom.com/wiki/Cave_Johnson",
		Summary:      "CEO of Aperture Science",
	}

	err := UpdateProfile(profile)
	if err != nil {
		t.Fatalf("Failed to update profile: %v", err)
	}

	// Get profile
	retrieved_profile, err := GetProfile(test_user.Id)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}

	// Test retrieved values
	if retrieved_profile.FirstName != profile.FirstName {
		t.Errorf("expected FirstName %s, got %s", profile.FirstName, retrieved_profile.FirstName)
	}

	if retrieved_profile.LastName != profile.LastName {
		t.Errorf("expected LastName %s, got %s", profile.LastName, retrieved_profile.LastName)
	}

	if retrieved_profile.Phone != profile.Phone {
		t.Errorf("expected Phone %s, got %s", profile.Phone, retrieved_profile.Phone)
	}

	if retrieved_profile.City != profile.City {
		t.Errorf("expected City %s, got %s", profile.City, retrieved_profile.City)
	}

	if retrieved_profile.State != profile.State {
		t.Errorf("expected State %s, got %s", profile.State, retrieved_profile.State)
	}

	if retrieved_profile.Country != profile.Country {
		t.Errorf("expected Country %s, got %s", profile.Country, retrieved_profile.Country)
	}

	if retrieved_profile.LinkedinURL != profile.LinkedinURL {
		t.Errorf("expected LinkedinURL %s, got %s", profile.LinkedinURL, retrieved_profile.LinkedinURL)
	}

	if retrieved_profile.PortfolioURL != profile.PortfolioURL {
		t.Errorf("expected PortfolioURL %s, got %s", profile.PortfolioURL, retrieved_profile.PortfolioURL)
	}

	if retrieved_profile.Summary != profile.Summary {
		t.Errorf("expected Summary %s, got %s", profile.Summary, retrieved_profile.Summary)
	}
}
