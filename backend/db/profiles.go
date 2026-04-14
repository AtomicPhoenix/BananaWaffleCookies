package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Profile struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	FirstName          string    `json:"first_name,omitempty"`
	LastName           string    `json:"last_name,omitempty"`
	Phone              string    `json:"phone,omitempty"`
	Location           string    `json:"location,omitempty"`
	City               string    `json:"city,omitempty"`
	State              string    `json:"state,omitempty"`
	Country            string    `json:"country,omitempty"`
	Headline           string    `json:"headline,omitempty"`
	LinkedinURL        string    `json:"linkedin_url,omitempty"`
	PortfolioURL       string    `json:"portfolio_url,omitempty"`
	Summary            string    `json:"summary,omitempty"`
	PreferredCity      string    `json:"preferred_city,omitempty"`
	PreferredState     string    `json:"preferred_state,omitempty"`
	PreferredRole      string    `json:"preferred_role,omitempty"`
	PreferredSalaryMin int       `json:"preferred_salary_min,omitempty"`
	PreferredSalaryMax int       `json:"preferred_salary_max,omitempty"`
	WorkMode           string    `json:"work_mode,omitempty"`
	CompletionPercent  int       `json:"completion_percent"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func createProfile(uid int) error {
	_, err := DbConn.Exec(context.Background(), "INSERT INTO profiles (user_id) VALUES ($1)", uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create profile: %v\n", err)
		return err
	}
	return nil
}

func UpdateProfile(profile Profile) error {
	query := `UPDATE profiles SET
		first_name = $1,
		last_name = $2,
		phone = $3,
		location = $4,
		city = $5,
		state = $6,
		country = $7,
		headline = $8,
		linkedin_url = $9,
		portfolio_url = $10,
		summary = $11,
		preferred_city = $12,
		preferred_state = $13,
		preferred_role = $14,
		preferred_salary_min = $15,
		preferred_salary_max = $16,
		work_mode = $17,
		completion_percent = $18,
		updated_at = NOW()
	WHERE user_id = $19`

	if profile.CompletionPercent == 0 {
		profile.SetCompletionPercent()
	}

	_, err := DbConn.Exec(
		context.Background(),
		query,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.Location,
		profile.City,
		profile.State,
		profile.Country,
		profile.Headline,
		profile.LinkedinURL,
		profile.PortfolioURL,
		profile.Summary,
		profile.PreferredCity,
		profile.PreferredState,
		profile.PreferredRole,
		profile.PreferredSalaryMin,
		profile.PreferredSalaryMax,
		profile.WorkMode,
		profile.CompletionPercent,
		profile.UserID,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update profile: %v\n", err)
		return err
	}
	return nil
}

func GetProfile(uid int) (Profile, error) {
	var profile Profile

	var firstName, lastName, phone, location, city, state, country sql.NullString
	var headline, linkedinURL, portfolioURL, summary sql.NullString
	var preferredCity, preferredState, preferredRole, workMode sql.NullString
	var preferredSalaryMin, preferredSalaryMax sql.NullInt64

	err := DbConn.QueryRow(context.Background(), `
		SELECT
			first_name,
			last_name,
			phone,
			location,
			city,
			state,
			country,
			headline,
			linkedin_url,
			portfolio_url,
			summary,
			preferred_city,
			preferred_state,
			preferred_role,
			preferred_salary_min,
			preferred_salary_max,
			work_mode,
			completion_percent
		FROM profiles
		WHERE user_id = $1
	`, uid).Scan(
		&firstName,
		&lastName,
		&phone,
		&location,
		&city,
		&state,
		&country,
		&headline,
		&linkedinURL,
		&portfolioURL,
		&summary,
		&preferredCity,
		&preferredState,
		&preferredRole,
		&preferredSalaryMin,
		&preferredSalaryMax,
		&workMode,
		&profile.CompletionPercent,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return Profile{}, err
	}

	profile.UserID = uid
	profile.FirstName = extractValue(firstName)
	profile.LastName = extractValue(lastName)
	profile.Phone = extractValue(phone)
	profile.Location = extractValue(location)
	profile.City = extractValue(city)
	profile.State = extractValue(state)
	profile.Country = extractValue(country)
	profile.Headline = extractValue(headline)
	profile.LinkedinURL = extractValue(linkedinURL)
	profile.PortfolioURL = extractValue(portfolioURL)
	profile.Summary = extractValue(summary)
	profile.PreferredCity = extractValue(preferredCity)
	profile.PreferredState = extractValue(preferredState)
	profile.PreferredRole = extractValue(preferredRole)
	profile.WorkMode = extractValue(workMode)

	if preferredSalaryMin.Valid {
		profile.PreferredSalaryMin = int(preferredSalaryMin.Int64)
	}
	if preferredSalaryMax.Valid {
		profile.PreferredSalaryMax = int(preferredSalaryMax.Int64)
	}

	return profile, nil
}

func extractValue(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}

func (profile *Profile) SetCompletionPercent() {
	profile.CompletionPercent = profile.getProfileCompletionPercent()
}

// Calculate what percent of the profile is filled out (returns an int from 0 to 100)
func (profile *Profile) getProfileCompletionPercent() int {
	var filledFields int = 0
	var numFields int = 17

	if profile.FirstName != "" {
		filledFields++
	}
	if profile.LastName != "" {
		filledFields++
	}
	if profile.Phone != "" {
		filledFields++
	}
	if profile.Location != "" {
		filledFields++
	}
	if profile.City != "" {
		filledFields++
	}
	if profile.State != "" {
		filledFields++
	}
	if profile.Country != "" {
		filledFields++
	}
	if profile.Headline != "" {
		filledFields++
	}
	if profile.LinkedinURL != "" {
		filledFields++
	}
	if profile.PortfolioURL != "" {
		filledFields++
	}
	if profile.Summary != "" {
		filledFields++
	}
	if profile.PreferredCity != "" {
		filledFields++
	}
	if profile.PreferredState != "" {
		filledFields++
	}
	if profile.PreferredRole != "" {
		filledFields++
	}
	if profile.PreferredSalaryMin > 0 {
		filledFields++
	}
	if profile.PreferredSalaryMax > 0 {
		filledFields++
	}
	if profile.WorkMode != "" {
		filledFields++
	}
	return int((float32(filledFields) / float32(numFields)) * 100)
}
