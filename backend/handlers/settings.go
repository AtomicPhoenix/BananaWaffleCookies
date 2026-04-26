package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
)

// Handler for /api/settings (PUT)
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userSettings := struct {
		email    string
		password string
	}{}

	err := json.NewDecoder(r.Body).Decode(&userSettings)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var tokenInfo Claim
	err, tokenInfo = GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusBadRequest)
		settings.Logger.Error("Failed to update settings; Failed to grab auth token information", "err", err)
		return

	}

	err = db.UpdateUserPassword(tokenInfo.Uid, userSettings.password)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusBadRequest)
		settings.Logger.Error("Failed to update settings; Failed to update password", "err", err)
		return
	}

	err = db.UpdateUserEmail(tokenInfo.Uid, userSettings.email)
	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusBadRequest)
		settings.Logger.Error("Failed to update settings; Failed to update email", "err", err)
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
		settings.Logger.Error("Failed to get settings; Failed to grab user token", "err", err)
		return

	}
	var uid int = tokenInfo.Uid

	user_email, err := db.GetUserEmail(uid)
	if err != nil {
		http.Error(w, "Failed to get settings", http.StatusBadRequest)
		settings.Logger.Error("Failed to get settings; Failed to grab user email", "err", err)
		return
	}

	settings := struct {
		Email string
	}{Email: user_email}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
