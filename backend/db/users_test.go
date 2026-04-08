package db

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

var test_uid int

func TestCreateUser(t *testing.T) {
	email := "test_user@example.com"
	password := "password123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := User{
		Email:         email,
		Password_hash: string(hash),
	}

	uid, err := RegisterUser(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	test_uid = uid
}

func TestUpdateUserEmail(t *testing.T) {
	new_email := "test_user_updated@example.com"
	err := UpdateUserEmail(test_uid, new_email)
	if err != nil {
		t.Fatalf("Failed to update user email: %v", err)
	}

	email, err := GetUserEmail(test_uid)
	if email != new_email || err != nil {
		t.Fatalf("Failed to update user email: %v", err)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	password := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		t.Fatalf("Failed to hash new password: %v", err)
	}

	err = UpdateUserPassword(test_uid, string(hash))
	if err != nil {
		t.Fatalf("Failed to update user password: %v", err)
	}

}

func TestGetUser(t *testing.T) {
	_, err := GetUserByID(test_uid)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	_, err := DbConn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", test_uid)
	if err != nil {
		t.Fatalf("failed to delete test user: %v", err)
	}
}
