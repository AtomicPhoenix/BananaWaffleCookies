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

func createProfile(uid int) bool {
	err := DbConn.QueryRow(context.Background(), "INSERT INTO profiles (user_id) VALUES ($1)", uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create profile: %v\n", err)
		return false
	}
	return true
}
