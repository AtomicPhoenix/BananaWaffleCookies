package db

import (
	"context"
	"fmt"
	"os"
)

type User struct {
	Id            int
	Email         string
	Password_hash string
}

func GetUser(email string) (User, error) {
	var user User
	err := DbConn.QueryRow(context.Background(), "SELECT id, email, password_hash FROM users WHERE email=$1", email).Scan(&user.Id, &user.Email, &user.Password_hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user: %v\n", err)
		return user, err
	}
	return user, nil
}

func RegisterUser(user User) (int, error) {
	var uid int
	fmt.Printf("Registering user %v\n", user)
	err := DbConn.QueryRow(context.Background(), "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", user.Email, user.Password_hash).Scan(&uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert user into database: %v\n", err)
		return uid, err
	}
	createProfile(uid)
	return uid, nil
}

func LoginUser(user User) (bool, int) {
	var uid int = -1
	err := DbConn.QueryRow(context.Background(), "SELECT uid FROM users WHERE email=$1 AND password_hash=$2", user.Email, user.Password_hash).Scan(&uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to authenticate user: %v\n", err)
		return false, -1
	}
	return true, uid
}
