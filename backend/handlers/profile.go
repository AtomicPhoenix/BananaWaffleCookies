package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
)

// Handler for /api/profile (PUT)
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusUnauthorized)
		settings.Logger.Error("Failed to update profile; Failed to grab auth token information", "err", err)
		return
	}

	var profile db.Profile

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&profile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		settings.Logger.Error("Failed to update profile; Failed to decode request body", "err", err)
		return
	}

	profile.UserID = tokenInfo.Uid
	profile.SetCompletionPercent()

	if err := db.UpdateProfile(profile); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update profile", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for /api/profile (GET)
func GetProfile(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusUnauthorized)
		settings.Logger.Error("Failed to get profile; Failed to grab auth token information", "err", err)
		return
	}

	profile, err := db.GetProfile(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Internal server error when getting profile", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get profile", "err", err)
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
		settings.Logger.Error("Failed to add education; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to add education; Failed to grab auth token information", "err", err)
		return
	}

	edu.UserID = tokenInfo.Uid

	id, err := db.InsertProfileEducation(edu)
	if err != nil {
		http.Error(w, "Failed to insert education", http.StatusBadRequest)
		settings.Logger.Error("Failed to insert education", "err", err)
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
		settings.Logger.Error("Failed to get education; Failed to grab auth token information", "err", err)
		return
	}

	edu, err := db.GetProfileEducation(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to fetch education", http.StatusBadRequest)
		settings.Logger.Error("Failed to fetch education", "err", err)
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
		settings.Logger.Error("Failed to delete education; Failed to grab auth token information", "err", err)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing experience id", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete education; Missing id")
		return
	}

	var eduID int
	fmt.Sscan(idParam, &eduID)

	if err := db.DeleteProfileEducation(tokenInfo.Uid, eduID); err != nil {
		http.Error(w, "Failed to delete education", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete education", "err", err)
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
		settings.Logger.Error("Failed to update education; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to update education; Failed to grab auth token information", "err", err)
		return
	}

	edu.UserID = tokenInfo.Uid

	if err := db.UpdateProfileEducation(edu); err != nil {
		http.Error(w, "Failed to update education", http.StatusBadRequest)
		settings.Logger.Error("Failed to update education", "err", err)
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
		settings.Logger.Error("Failed to reorder education; Failed to grab auth token information", "err", err)
		return
	}

	reordering := struct {
		ID       int `json:"id"`
		Position int `json:"position"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&reordering); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		settings.Logger.Error("Failed to reorder education; Failed to decode request body", "err", err)
		return
	}

	if err = db.ReorderProfileEducation(tokenInfo.Uid, reordering.ID, reordering.Position); err != nil {
		http.Error(w, "Failed to reorder Education", http.StatusBadRequest)
		settings.Logger.Error("Failed to reorder education", "err", err)
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
		settings.Logger.Error("Failed to add experience; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to add experience; Failed to grab auth token information", "err", err)
		return
	}

	exp.UserID = tokenInfo.Uid

	id, err := db.InsertProfileExperience(exp)
	if err != nil {
		http.Error(w, "Failed to insert experience", http.StatusBadRequest)
		settings.Logger.Error("Failed to insert experience", "err", err)
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
		settings.Logger.Error("Failed to get experiences; Failed to grab auth token information", "err", err)
		return
	}

	exps, err := db.GetProfileExperiences(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to fetch experiences", http.StatusBadRequest)
		settings.Logger.Error("Failed to fetch experiences", "err", err)
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
		settings.Logger.Error("Failed to delete experience; Failed to grab auth token information", "err", err)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing experience id", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete experience; Missing id")
		return
	}

	var expID int
	fmt.Sscan(idParam, &expID)

	if err := db.DeleteProfileExperience(tokenInfo.Uid, expID); err != nil {
		http.Error(w, "Failed to delete experience", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete experience", "err", err)
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
		settings.Logger.Error("Failed to update experience; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to update experience; Failed to grab auth token information", "err", err)
		return
	}

	exp.UserID = tokenInfo.Uid

	if err := db.UpdateProfileExperience(exp); err != nil {
		http.Error(w, "Failed to update experience", http.StatusBadRequest)
		settings.Logger.Error("Failed to update experience", "err", err)
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
		settings.Logger.Error("Failed to reorder experiences; Failed to grab auth token information", "err", err)
		return
	}

	reordering := struct {
		ID       int `json:"id"`
		Position int `json:"position"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&reordering); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		settings.Logger.Error("Failed to reorder experiences; Failed to decode request body", "err", err)
		return
	}

	if err = db.ReorderProfileExperience(tokenInfo.Uid, reordering.ID, reordering.Position); err != nil {
		http.Error(w, "Failed to reorder Experience", http.StatusBadRequest)
		settings.Logger.Error("Failed to reorder experiences", "err", err)
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
		settings.Logger.Error("Failed to add skill; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to add skill; Failed to grab auth token information", "err", err)
		return
	}

	skill.UserID = tokenInfo.Uid

	id, err := db.InsertProfileSkill(skill)
	if err != nil {
		http.Error(w, "Failed to insert skill", http.StatusBadRequest)
		settings.Logger.Error("Failed to insert skill", "err", err)
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
		settings.Logger.Error("Failed to get skills; Failed to grab auth token information", "err", err)
		return
	}

	skills, err := db.GetProfileSkills(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to fetch skills", http.StatusBadRequest)
		settings.Logger.Error("Failed to fetch skills", "err", err)
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
		settings.Logger.Error("Failed to delete skill; Failed to grab auth token information", "err", err)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete skill; Missing id")
		return
	}

	var skillID int
	fmt.Sscan(idParam, &skillID)

	if err := db.DeleteProfileSkill(tokenInfo.Uid, skillID); err != nil {
		http.Error(w, "Failed to delete skill", http.StatusBadRequest)
		settings.Logger.Error("Failed to delete skill", "err", err)
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
		settings.Logger.Error("Failed to update skill; Failed to decode request body", "err", err)
		return
	}

	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to update skill; Failed to grab auth token information", "err", err)
		return
	}

	skill.UserID = tokenInfo.Uid

	if err := db.UpdateProfileSkill(skill); err != nil {
		http.Error(w, "Failed to update skill", http.StatusBadRequest)
		settings.Logger.Error("Failed to update skill", "err", err)
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
		settings.Logger.Error("Failed to reorder skill; Failed to grab auth token information", "err", err)
		return
	}

	reordering := struct {
		ID       int `json:"id"`
		Position int `json:"position"`
	}{}

	if err = json.NewDecoder(r.Body).Decode(&reordering); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		settings.Logger.Error("Failed to reorder skill; Failed to decode request body", "err", err)
		return
	}

	if err = db.ReorderProfileSkill(tokenInfo.Uid, reordering.ID, reordering.Position); err != nil {
		http.Error(w, "Failed to reorder skill", http.StatusBadRequest)
		settings.Logger.Error("Failed to reorder skill", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Skill reordered successfully "}`)
}

// Handler for PUT /api/profile/preferences
func UpdateProfilePreferences(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to update preferences; Failed to grab auth token information", "err", err)
		return
	}

	var req db.Profile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		settings.Logger.Error("Failed to update preferences; Failed to decode request body", "err", err)
		return
	}

	if err := db.UpdateProfilePreferences(tokenInfo.Uid, req); err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		settings.Logger.Error("Failed to update preferences", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler for GET /api/profile/preferences
func GetProfilePreferences(w http.ResponseWriter, r *http.Request) {
	err, tokenInfo := GrabToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		settings.Logger.Error("Failed to get preferences; Failed to grab auth token information", "err", err)
		return
	}

	prefs, err := db.GetProfilePreferences(tokenInfo.Uid)
	if err != nil {
		http.Error(w, "Failed to get preferences", http.StatusInternalServerError)
		settings.Logger.Error("Failed to get preferences", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prefs)
}
