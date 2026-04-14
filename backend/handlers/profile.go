package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bananawafflecookies.com/m/v2/db"
)

// Handler for /api/profile (PUT)
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusUnauthorized)
		fmt.Fprintf(os.Stderr, "Failed to update profile: %v\n", err)
		return

	}

	var profile db.Profile

	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile.UserID = tokenInfo.Uid
	profile.SetCompletionPercent()

	err = db.UpdateProfile(profile)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Internal server error when updating profile: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for /api/profile (GET)
func GetProfile(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusUnauthorized)
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return

	}
	var uid int = tokenInfo.Uid

	profile, err := db.GetProfile(uid)
	if err != nil {
		http.Error(w, "Internal server error when getting profile", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
