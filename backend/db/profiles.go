package db

import (
	"context"
	"fmt"
	"os"
)

type Profile struct {
	user_id       int
	first_name    string
	last_name     string
	phone         string
	city          string
	state         string
	country       string
	linkedin_url  string
	portfolio_url string
	summary       string
	completion    int
}

func createProfile(profile Profile) bool {
	var sql string = `INSERT INTO 
				profiles (user_id, first_name, last_name, phone, city, state, country, linkedin_url, portfolio_url, summary, completion_percent) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	err := DbConn.QueryRow(context.Background(), sql, profile.user_id, profile.first_name, profile.last_name, profile.phone, profile.city, profile.state, profile.country, profile.linkedin_url, profile.portfolio_url, profile.summary, profile.completion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create profile: %v\n", err)
		return false
	}
	return true
}
