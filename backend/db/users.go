package db

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            int
	Email         string
	Password_hash string
}

func GetUser(email string) (User, error) {
	var user User
	err := DbConn.QueryRow(context.Background(),
		"SELECT id, email, password_hash FROM users WHERE email=$1",
		email,
	).Scan(&user.Id, &user.Email, &user.Password_hash)

	if err != nil {
		return user, fmt.Errorf("Failed to get user by email=%s: %w", email, err)
	}
	return user, nil
}

func GetUserByID(uid int) (User, error) {
	var user User
	err := DbConn.QueryRow(context.Background(),
		"SELECT id, email, password_hash FROM users WHERE id=$1",
		uid,
	).Scan(&user.Id, &user.Email, &user.Password_hash)

	if err != nil {
		return user, fmt.Errorf("Failed to get user by id=%d: %w", uid, err)
	}
	return user, nil
}

func RegisterUser(user User) (int, error) {
	var uid int
	err := DbConn.QueryRow(context.Background(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		user.Email,
		user.Password_hash,
	).Scan(&uid)

	if err != nil {
		return uid, fmt.Errorf("Failed to insert user (email=%s): %w", user.Email, err)
	}

	err = createProfile(uid)
	if err != nil {
		return uid, fmt.Errorf("Failed to create profile for user_id=%d: %w", uid, err)
	}

	return uid, nil
}

func UpdateUserPassword(uid int, new_password string) error {
	password_bytes, err := bcrypt.GenerateFromPassword([]byte(new_password), 12)
	if err != nil {
		return fmt.Errorf("Failed to generate password hash for user_id=%d: %w", uid, err)
	}

	password_hash := string(password_bytes)

	_, err = DbConn.Exec(context.Background(),
		"UPDATE users SET password_hash=$1 WHERE id=$2",
		password_hash,
		uid,
	)

	if err != nil {
		return fmt.Errorf("Failed to update password for user_id=%d: %w", uid, err)
	}

	return nil
}

func UpdateUserEmail(uid int, new_email string) error {
	_, err := DbConn.Exec(context.Background(),
		"UPDATE users SET email=$1 WHERE id=$2",
		new_email,
		uid,
	)

	if err != nil {
		return fmt.Errorf("Failed to update email for user_id=%d: %w", uid, err)
	}

	return nil
}

func GetUserEmail(uid int) (string, error) {
	var email string
	err := DbConn.QueryRow(context.Background(),
		"SELECT email FROM users WHERE id=$1",
		uid,
	).Scan(&email)

	if err != nil {
		return "", fmt.Errorf("Failed to get email for user_id=%d: %w", uid, err)
	}

	return email, nil
}
