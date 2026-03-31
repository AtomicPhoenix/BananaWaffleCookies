package db

import (
	"context"
	"fmt"
	"os"
)

type Profile struct {
	user_id            int
	first_name         string
	last_name          string
	phone              string
	city               string
	state              string
	country            string
	linkedin_url       string
	portfolio_url      string
	summary            string
	completion_percent int
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
	_, err := DbConn.Exec(context.Background(), sql, profile.first_name, profile.last_name, profile.phone, profile.city, profile.state, profile.country, profile.linkedin_url, profile.portfolio_url, profile.summary, profile.completion_percent, profile.user_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert user into database: %v\n", err)
		return err
	}
	return nil
}
