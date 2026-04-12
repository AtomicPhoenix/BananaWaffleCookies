package db

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func createTestUser(t *testing.T) User {
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

	user.Id = uid
	return user
}

func deleteTestUser(t *testing.T, test_uid int) {
	_, err := DbConn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", test_uid)
	if err != nil {
		t.Fatalf("failed to delete test user: %v", err)
	}
}

func TestCreateAndDeleteUser(t *testing.T) {
	user := createTestUser(t)
	deleteTestUser(t, user.Id)
}

func TestUpdateUserEmail(t *testing.T) {
	user := createTestUser(t)
	defer deleteTestUser(t, user.Id)

	new_email := "test_user_updated@example.com"
	err := UpdateUserEmail(user.Id, new_email)
	if err != nil {
		t.Fatalf("Failed to update user email: %v", err)
	}

	email, err := GetUserEmail(user.Id)
	if email != new_email || err != nil {
		t.Fatalf("Failed to update user email: %v", err)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	user := createTestUser(t)
	defer deleteTestUser(t, user.Id)

	password := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		t.Fatalf("Failed to hash new password: %v", err)
	}

	err = UpdateUserPassword(user.Id, string(hash))
	if err != nil {
		t.Fatalf("Failed to update user password: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	user := createTestUser(t)
	defer deleteTestUser(t, user.Id)

	_, err := GetUserByID(user.Id)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
}
