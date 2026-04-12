package db

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
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

func GetUserByID(uid int) (User, error) {
	var user User
	err := DbConn.QueryRow(context.Background(), "SELECT id, email, password_hash FROM users WHERE id=$1", uid).Scan(&user.Id, &user.Email, &user.Password_hash)
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
	err = createProfile(uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create user profile: %v\n", err)
		return uid, err
	}
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

func UpdateUserPassword(uid int, new_password string) error {
	password_bytes, err := bcrypt.GenerateFromPassword([]byte(new_password), 12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate password hash: %v\n", err)
		return err
	}
	password_hash := string(password_bytes)
	_, err = DbConn.Exec(context.Background(), "UPDATE users SET password_hash=$1 WHERE id=$2", password_hash, uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update user password: %v\n", err)
		return err
	}
	return nil
}

func UpdateUserEmail(uid int, new_email string) error {
	_, err := DbConn.Exec(context.Background(), "UPDATE users SET email=$1 WHERE id=$2", new_email, uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update user email: %v\n", err)
		return err
	}
	return nil
}

func GetUserEmail(uid int) (string, error) {
	var email string
	err := DbConn.QueryRow(context.Background(), "SELECT email FROM users WHERE id=$1", uid).Scan(&email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update get email: %v\n", err)
		return "", err
	}
	return email, nil
}
