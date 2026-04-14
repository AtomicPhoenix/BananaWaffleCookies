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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&profile)
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

// Handler for /api/profile/education (POST)
func AddProfileEducation(w http.ResponseWriter, r *http.Request) {
	var edu db.ProfileEducation

	if err := json.NewDecoder(r.Body).Decode(&edu); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	edu.UserID = tokenInfo.Uid

	id, err := db.InsertProfileEducation(edu)
	if err != nil {
		http.Error(w, "Failed to insert education", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to insert education: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Handler for /api/profile/education (GET)
func GetProfileEducation(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	edu, err := db.GetProfileEducation(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to fetch education", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to fetch education: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(edu)
}

// Handler for /api/profile/education/{id} (DELETE)
func DeleteProfileEducation(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing experience id", http.StatusBadRequest)
		return
	}

	var eduID int
	fmt.Sscan(idParam, &eduID)

	err = db.DeleteProfileEducation(tokenInfo.Uid, eduID)
	if err != nil {
		http.Error(w, "Failed to delete education", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to delete education: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Education deleted successfully "}`)
}

// Handler for /api/profile/education (PUT)
func UpdateProfileEducation(w http.ResponseWriter, r *http.Request) {
	var edu db.ProfileEducation

	if err := json.NewDecoder(r.Body).Decode(&edu); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	edu.UserID = tokenInfo.Uid

	err = db.UpdateProfileEducation(edu)
	if err != nil {
		http.Error(w, "Failed to update education", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update education: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Education updated successfully"}`)
}

// Handler for /api/profile/education/reorder
func ReorderProfileEducation(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reordering := struct {
		ID       int `json:"id"`
		Position int `json:"position"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&reordering); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.ReorderProfileEducation(tokenInfo.Uid, reordering.ID, reordering.Position)
	if err != nil {
		http.Error(w, "Failed to reorder Education", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to reorder Education: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Education reordered successfully "}`)
}

// Handler for /api/profile/experiences (POST)
func AddProfileExperience(w http.ResponseWriter, r *http.Request) {
	var exp db.ProfileExperiences

	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exp.UserID = tokenInfo.Uid

	id, err := db.InsertProfileExperience(exp)
	if err != nil {
		http.Error(w, "Failed to insert experience", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to insert experience: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Handler for /api/profile/experiences (GET)
func GetProfileExperiences(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exps, err := db.GetProfileExperiences(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to fetch experiences", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to fetch experiences: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exps)
}

// Handler for /api/profile/experiences/{id} (DELETE)
func DeleteProfileExperience(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing experience id", http.StatusBadRequest)
		return
	}

	var expID int
	fmt.Sscan(idParam, &expID)

	err = db.DeleteProfileExperience(tokenInfo.Uid, expID)
	if err != nil {
		http.Error(w, "Failed to delete experience", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to delete experience: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Experience deleted successfully "}`)
}

// Handler for /api/profile/experiences (PUT)
func UpdateProfileExperience(w http.ResponseWriter, r *http.Request) {
	var exp db.ProfileExperiences

	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exp.UserID = tokenInfo.Uid

	err = db.UpdateProfileExperience(exp)
	if err != nil {
		http.Error(w, "Failed to update experience", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update experience: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Experience updated successfully"}`)
}

// Handler for /api/profile/experiences/reorder
func ReorderProfileExperiences(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reordering := struct {
		ID       int `json:"id"`
		Position int `json:"position"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&reordering); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.ReorderProfileExperience(tokenInfo.Uid, reordering.ID, reordering.Position)
	if err != nil {
		http.Error(w, "Failed to reorder Experience", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to reorder Experience: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Experience reordered successfully "}`)
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

// Handler for /api/profile/skills (PUT)
func UpdateProfileSkill(w http.ResponseWriter, r *http.Request) {
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

	err = db.UpdateProfileSkill(skill)
	if err != nil {
		http.Error(w, "Failed to update skill", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to update skill: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Skill updated successfully"}`)
}

// Handler for /api/profile/skills/reorder
func ReorderProfileSkill(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reordering := struct {
		ID       int `json:"id"`
		Position int `json:"position"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&reordering); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.ReorderProfileSkill(tokenInfo.Uid, reordering.ID, reordering.Position)
	if err != nil {
		http.Error(w, "Failed to reorder skill", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Failed to reorder skill: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Skill reordered successfully "}`)
}

// Handler for PUT /api/profile/preferences
func UpdateProfilePreferences(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req db.Profile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = db.UpdateProfilePreferences(tokenInfo.Uid, req)
	if err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for GET /api/profile/preferences
func GetProfilePreferences(w http.ResponseWriter, r *http.Request) {
	var tokenInfo Claim
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	prefs, err := db.GetProfilePreferences(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prefs)
}
