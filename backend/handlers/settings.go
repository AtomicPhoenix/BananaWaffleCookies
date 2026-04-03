package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bananawafflecookies.com/m/v2/db"
)

// Handler for /api/settings (PUT)
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	settings := struct {
		email    string
		password string
	}{}

	err := json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var tokenInfo Claim
	err, tokenInfo = GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update settings; Failed to grab auth token information: %v\n", err)
		return

	}

	err = db.UpdateUserPassword(tokenInfo.Uid, settings.password)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update settings; Failed to update password: %v\n", err)
		return
	}

	err = db.UpdateUserEmail(tokenInfo.Uid, settings.email)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update settings; Failed to update email: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Settings successfully updated."}`)
}

// Handler for /api/settings (GET)
func GetSettings(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get settings", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get settings: %v\n", err)
		return

	}
	var uid int = tokenInfo.Uid

	user_email, err := db.GetUserEmail(uid)
	if err != nil {
		http.Error(w, "Failed to get settings", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get settings: %v\n", err)
		return
	}

	settings := struct {
		Email string
	}{Email: user_email}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
