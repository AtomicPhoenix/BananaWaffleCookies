package db

import (
	"context"
	"fmt"
	"os"
)

type User struct {
	email         string
	password_hash string
}

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

func GetUserID(email string) int {
	var id int
	err := DbConn.QueryRow(context.Background(), "SELECT id FROM users WHERE email=$1", email).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user: %v\n", err)
	}
	return id
}

func RegisterUser(user User, profile Profile) bool {
	err := DbConn.QueryRow(context.Background(), "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", user.email, user.password_hash).Scan(&profile.user_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return false
	}
	return createProfile(profile)
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

func LoginUser(email string, password_hash string) (bool, int) {
	var uid int = -1
	err := DbConn.QueryRow(context.Background(), "SELECT uid FROM users WHERE email=$1 AND password_hash=$2", email, password_hash).Scan(&uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to authenticate user: %v\n", err)
		return false, -1
	}
	return true, uid
}
