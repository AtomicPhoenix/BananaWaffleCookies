package handlers

import (
	"context"
	"testing"

	"bananawafflecookies.com/m/v2/db"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func createTestUser(t *testing.T) db.User {
	email := "test_user@example.com"
	password := "password123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := db.User{
		Email:         email,
		Password_hash: string(hash),
	}

	uid, err := db.RegisterUser(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	user.Id = uid
	return user
}

func deleteTestUser(t *testing.T, test_uid int) {
	_, err := db.DbConn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", test_uid)
	if err != nil {
		t.Fatalf("failed to delete test user: %v", err)
	}
}

// Setup database connection for testing
func TestMain(t *testing.T) {
	godotenv.Load("../.env")
	err := db.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v\n", err)
	}
	initAuth()
}
