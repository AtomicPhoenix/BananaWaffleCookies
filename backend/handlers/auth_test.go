package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

// Helper function to generate auth token for protected paths
func getAuthCookie(t *testing.T, uid int, email string) *http.Cookie {
	claims := map[string]interface{}{
		"id":    fmt.Sprintf("%v", uid),
		"email": email,
	}
	_, tokenString, err := AuthToken.Encode(claims)
	if err != nil {
		t.Errorf("JWT encode error: %v", err)
	}
	return &http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	}
}
