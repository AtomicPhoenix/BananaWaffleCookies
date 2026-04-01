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
	var profile db.Profile

	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var tokenInfo Claim
	err, tokenInfo = GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update profile: %v\n", err)
		return

	}
	profile.UserID = tokenInfo.Uid
	profile.SetCompletionPercent()

	err = db.UpdateProfile(profile)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update profile: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Profile successfully updated."}`)
}
