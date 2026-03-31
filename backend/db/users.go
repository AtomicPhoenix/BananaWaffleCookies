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
		fmt.Fprintf(os.Stderr, "Failed to insert user into database: %v\n", err)
		return false
	}
	return createProfile(profile)
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
