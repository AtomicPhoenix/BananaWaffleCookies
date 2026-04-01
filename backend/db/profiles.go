package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Profile struct {
	Id                int       `json:"id"`
	UserID            int       `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Phone             string    `json:"phone"`
	City              string    `json:"city"`
	State             string    `json:"state"`
	Country           string    `json:"country"`
	LinkedinURL       string    `json:"linkedin_url"`
	PortfolioURL      string    `json:"portfolio_url"`
	Summary           string    `json:"summary"`
	CompletionPercent int       `json:"completion_percent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func createProfile(uid int) bool {
	_, err := DbConn.Exec(context.Background(), "INSERT INTO profiles (user_id) VALUES ($1)", uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create profile: %v\n", err)
		return false
	}
	return true
}

func UpdateProfile(profile Profile) error {
	var sql string = `UPDATE profiles
				SET first_name = $1, last_name = $2, phone = $3, city = $4, state = $5, country = $6, linkedin_url = $7, portfolio_url = $8, summary = $9, completion_percent = $10
				WHERE user_id = $11`
	_, err := DbConn.Exec(context.Background(), sql, profile.FirstName, profile.LastName, profile.Phone, profile.City, profile.State, profile.Country, profile.LinkedinURL, profile.PortfolioURL, profile.Summary, profile.CompletionPercent, profile.UserID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert user into database: %v\n", err)
		return err
	}
	return nil
}

func GetProfile(uid int) (Profile, error) {
	var profile Profile
	var first_name, last_name, phone, city, state, country, linkedin_url, portfolio, summary sql.NullString
	err := DbConn.QueryRow(context.Background(), `SELECT first_name, last_name, phone, city, state, country, linkedin_url, portfolio_url, summary, completion_percent FROM profiles WHERE user_id = $1`, uid).Scan(&first_name, &last_name, &phone, &city, &state, &country, &linkedin_url, &portfolio, &summary, &profile.CompletionPercent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to grab user from database: %v\n", err)
		return Profile{}, err
	}

	profile.FirstName = extractValue(first_name)
	profile.LastName = extractValue(last_name)
	profile.City = extractValue(city)
	profile.Phone = extractValue(phone)
	profile.State = extractValue(state)
	profile.Country = extractValue(country)
	profile.LinkedinURL = extractValue(linkedin_url)
	profile.PortfolioURL = extractValue(portfolio)
	profile.Summary = extractValue(summary)

	return profile, nil
}

func extractValue(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}

// Calculate what percent of the profile is filled out
func (profile *Profile) SetCompletionPercent() {
	var filledFields int = 0
	var numFields int = 9
	if profile.FirstName != "" {
		filledFields++
	}

	if profile.LastName != "" {
		filledFields++
	}

	if profile.Phone != "" {
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

	if profile.LinkedinURL != "" {
		filledFields++
	}

	if profile.PortfolioURL != "" {
		filledFields++
	}

	if profile.Summary != "" {
		filledFields++
	}

	profile.CompletionPercent = int(float32(filledFields/numFields) * 100)
}
