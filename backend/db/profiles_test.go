package db

import "testing"

func TestUpdateProfile(t *testing.T) {
	// Creating a test user also creates a profile automatically
	test_user := createTestUser(t)
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	profile := Profile{
		UserID:             test_user.Id,
		FirstName:          "Cave",
		LastName:           "Johnson",
		Phone:              "5132737724",
		Location:           "Somewhere in Michigan",
		City:               "Unknown",
		State:              "Unknown",
		Country:            "USA",
		Headline:           "CEO of Aperture Science",
		LinkedinURL:        "https://half-life.fandom.com/wiki/Cave_Johnson",
		PortfolioURL:       "https://half-life.fandom.com/wiki/Cave_Johnson",
		Summary:            "We do what we must because we can.",
		PreferredCity:      "Detroit",
		PreferredState:     "MI",
		PreferredRole:      "Scientist",
		PreferredSalaryMin: 100000000,
		PreferredSalaryMax: 200000000,
		WorkMode:           "onsite",
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

	// Test retrieved_profile values
	if retrieved_profile.FirstName != profile.FirstName {
		t.Errorf("Expected FirstName %s, got %s", profile.FirstName, retrieved_profile.FirstName)
	}
	if retrieved_profile.LastName != profile.LastName {
		t.Errorf("Expected LastName %s, got %s", profile.LastName, retrieved_profile.LastName)
	}
	if retrieved_profile.Phone != profile.Phone {
		t.Errorf("Expected Phone %s, got %s", profile.Phone, retrieved_profile.Phone)
	}
	if retrieved_profile.Location != profile.Location {
		t.Errorf("Expected Location %s, got %s", profile.Location, retrieved_profile.Location)
	}
	if retrieved_profile.City != profile.City {
		t.Errorf("Expected City %s, got %s", profile.City, retrieved_profile.City)
	}
	if retrieved_profile.State != profile.State {
		t.Errorf("Expected State %s, got %s", profile.State, retrieved_profile.State)
	}
	if retrieved_profile.Country != profile.Country {
		t.Errorf("Expected Country %s, got %s", profile.Country, retrieved_profile.Country)
	}

	// Professional info
	if retrieved_profile.Headline != profile.Headline {
		t.Errorf("Expected Headline %s, got %s", profile.Headline, retrieved_profile.Headline)
	}
	if retrieved_profile.LinkedinURL != profile.LinkedinURL {
		t.Errorf("Expected LinkedinURL %s, got %s", profile.LinkedinURL, retrieved_profile.LinkedinURL)
	}
	if retrieved_profile.PortfolioURL != profile.PortfolioURL {
		t.Errorf("Expected PortfolioURL %s, got %s", profile.PortfolioURL, retrieved_profile.PortfolioURL)
	}
	if retrieved_profile.Summary != profile.Summary {
		t.Errorf("Expected Summary %s, got %s", profile.Summary, retrieved_profile.Summary)
	}

	// Preferences
	if retrieved_profile.PreferredCity != profile.PreferredCity {
		t.Errorf("Expected PreferredCity %s, got %s", profile.PreferredCity, retrieved_profile.PreferredCity)
	}
	if retrieved_profile.PreferredState != profile.PreferredState {
		t.Errorf("Expected PreferredState %s, got %s", profile.PreferredState, retrieved_profile.PreferredState)
	}
	if retrieved_profile.PreferredRole != profile.PreferredRole {
		t.Errorf("Expected PreferredRole %s, got %s", profile.PreferredRole, retrieved_profile.PreferredRole)
	}
	if retrieved_profile.PreferredSalaryMin != profile.PreferredSalaryMin {
		t.Errorf("Expected PreferredSalaryMin %d, got %d", profile.PreferredSalaryMin, retrieved_profile.PreferredSalaryMin)
	}
	if retrieved_profile.PreferredSalaryMax != profile.PreferredSalaryMax {
		t.Errorf("Expected PreferredSalaryMax %d, got %d", profile.PreferredSalaryMax, retrieved_profile.PreferredSalaryMax)
	}
	if retrieved_profile.WorkMode != profile.WorkMode {
		t.Errorf("Expected WorkMode %s, got %s", profile.WorkMode, retrieved_profile.WorkMode)
	}

	// Confirm completion percent updated
	if retrieved_profile.CompletionPercent == profile.CompletionPercent {
		t.Errorf("Expected updated CompletionPercent %d, got %d", retrieved_profile.getProfileCompletionPercent(), retrieved_profile.CompletionPercent)
	}

	// Confirm completion percent check calculated correctly
	if retrieved_profile.CompletionPercent != retrieved_profile.getProfileCompletionPercent() {
		t.Errorf("Expected CompletionPercent %d, got %d", retrieved_profile.getProfileCompletionPercent(), retrieved_profile.CompletionPercent)
	}
}
