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

// Handler for /api/profile (GET)
func GetProfile(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return

	}
	var uid int = tokenInfo.Uid

	profile, err := db.GetProfile(uid)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// Handler for /api/profile/skills (POST)
func AddProfileSkill(w http.ResponseWriter, r *http.Request) {
	var skill db.ProfileSkills

	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	skill.UserID = tokenInfo.Uid

	id, err := db.InsertProfileSkill(skill)
	if err != nil {
		http.Error(w, "Failed to insert skill", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to insert skill: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Handler for /api/profile/skills (GET)
func GetProfileSkills(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	skills, err := db.GetProfileSkills(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to fetch skills", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to fetch skills: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

// Handler for /api/profile/skills/{id} (DELETE)
func DeleteProfileSkill(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	var skillID int
	fmt.Sscan(idParam, &skillID)

	err = db.DeleteProfileSkill(tokenInfo.Uid, skillID)
	if err != nil {
		http.Error(w, "Failed to delete skill", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to delete skill: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Skill deleted successfully "}`)
}
