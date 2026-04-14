package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Profile struct {
	Id                int       `json:"id"`
	UserID            int       `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Phone             string    `json:"phone"`
	City              string    `json:"city"`
	State             string    `json:"state"`
	Country           string    `json:"country"`
	LinkedinURL       string    `json:"linkedin_url"`
	PortfolioURL      string    `json:"portfolio_url"`
	Summary           string    `json:"summary"`
	CompletionPercent int       `json:"completion_percent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ProfileEducation struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Institution  string    `json:"institution"`
	Degree       string    `json:"degree,omitempty"`
	FieldOfStudy string    `json:"field_of_study,omitempty"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IsCurrent    bool      `json:"is_current"`
	Honors       string    `json:"honors,omitempty"`
	Gpa          float64   `json:"gpa,omitempty"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProfileExperiences struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ExperienceType string    `json:"experience_type"`
	Title          string    `json:"title"`
	Organization   string    `json:"organization,omitempty"`
	LocationText   string    `json:"location_text,omitempty"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	IsCurrent      bool      `json:"is_current"`
	Description    string    `json:"description,omitempty"`
	SortOrder      int       `json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProfileSkills struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	SkillName        string    `json:"skill_name"`
	Category         string    `json:"category,omitempty"`
	ProficiencyLabel string    `json:"proficiency_label,omitempty"`
	SortOrder        int       `json:"sort_order"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func createProfile(uid int) error {
	_, err := DbConn.Exec(context.Background(), "INSERT INTO profiles (user_id) VALUES ($1)", uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create profile: %v\n", err)
		return err
	}
	return nil
}

func UpdateProfile(profile Profile) error {
	var sql string = `UPDATE profiles
				SET first_name = $1, last_name = $2, phone = $3, city = $4, state = $5, country = $6, linkedin_url = $7, portfolio_url = $8, summary = $9, completion_percent = $10
				WHERE user_id = $11`
	_, err := DbConn.Exec(context.Background(), sql, profile.FirstName, profile.LastName, profile.Phone, profile.City, profile.State, profile.Country, profile.LinkedinURL, profile.PortfolioURL, profile.Summary, profile.CompletionPercent, profile.UserID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert user into database: %v\n", err)
		return err
	}
	return nil
}

func GetProfile(uid int) (Profile, error) {
	var profile Profile
	var first_name, last_name, phone, city, state, country, linkedin_url, portfolio, summary sql.NullString
	err := DbConn.QueryRow(context.Background(), `SELECT first_name, last_name, phone, city, state, country, linkedin_url, portfolio_url, summary, completion_percent FROM profiles WHERE user_id = $1`, uid).Scan(&first_name, &last_name, &phone, &city, &state, &country, &linkedin_url, &portfolio, &summary, &profile.CompletionPercent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to grab user from database: %v\n", err)
		return Profile{}, err
	}

	profile.Id = uid
	profile.FirstName = extractValue(first_name)
	profile.LastName = extractValue(last_name)
	profile.City = extractValue(city)
	profile.Phone = extractValue(phone)
	profile.State = extractValue(state)
	profile.Country = extractValue(country)
	profile.LinkedinURL = extractValue(linkedin_url)
	profile.PortfolioURL = extractValue(portfolio)
	profile.Summary = extractValue(summary)

	return profile, nil
}

func extractValue(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}

// Calculate what percent of the profile is filled out
func (profile *Profile) SetCompletionPercent() {
	var filledFields int = 0
	var numFields int = 9
	if profile.FirstName != "" {
		filledFields++
	}

	if profile.LastName != "" {
		filledFields++
	}

	if profile.Phone != "" {
		filledFields++
	}

	if profile.City != "" {
		filledFields++
	}

	if profile.State != "" {
		filledFields++
	}

	if profile.Country != "" {
		filledFields++
	}

	if profile.LinkedinURL != "" {
		filledFields++
	}

	if profile.PortfolioURL != "" {
		filledFields++
	}

	if profile.Summary != "" {
		filledFields++
	}

	profile.CompletionPercent = int(float32(filledFields/numFields) * 100)
}

func InsertProfileEducation(e ProfileEducation) (int, error) {
	var id int
	query := `
		INSERT INTO profile_education 
		(user_id, institution, degree, field_of_study, start_date, end_date, is_current, honors, gpa, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`

	err := DbConn.QueryRow(
		context.Background(),
		query,
		e.UserID,
		e.Institution,
		e.Degree,
		e.FieldOfStudy,
		e.StartDate,
		e.EndDate,
		e.IsCurrent,
		e.Honors,
		e.Gpa,
		e.SortOrder,
	).Scan(&id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert education into profile: %v", err)
		return 0, err
	}
	return id, nil
}

func GetProfileEducation(userID int) ([]ProfileEducation, error) {
	rows, err := DbConn.Query(context.Background(), `
		SELECT id, user_id, institution, degree, field_of_study,
		       start_date, end_date, is_current, honors, gpa,
		       sort_order, created_at, updated_at
		FROM profile_education
		WHERE user_id = $1
		ORDER BY sort_order ASC, id ASC
	`, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to grab education from profile: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []ProfileEducation
	for rows.Next() {
		var e ProfileEducation
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.Institution,
			&e.Degree,
			&e.FieldOfStudy,
			&e.StartDate,
			&e.EndDate,
			&e.IsCurrent,
			&e.Honors,
			&e.Gpa,
			&e.SortOrder,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to grab education from profile: %v", err)
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func DeleteProfileEducation(userID, educationID int) error {
	tag, err := DbConn.Exec(context.Background(), `
		DELETE FROM profile_education
		WHERE id = $1 AND user_id = $2
	`, educationID, userID)

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile education deletion")
	}

	return err
}

func UpdateProfileEducation(edu ProfileEducation) error {
	query := `
		UPDATE profile_education
		SET institution = $1,
		    degree = $2,
		    field_of_study = $3,
		    start_date = $4,
		    end_date = $5,
		    is_current = $6,
		    honors = $7,
		    gpa = $8,
		    sort_order = $9,
		    updated_at = NOW()
		WHERE id = $10 AND user_id = $11
	`

	tag, err := DbConn.Exec(
		context.Background(),
		query,
		edu.Institution,
		edu.Degree,
		edu.FieldOfStudy,
		edu.StartDate,
		edu.EndDate,
		edu.IsCurrent,
		edu.Honors,
		edu.Gpa,
		edu.SortOrder,
		edu.ID,
		edu.UserID,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update education: %v\n", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile education update")
	}

	return nil
}

func ReorderProfileEducation(userID int, eduID int, sortOrder int) error {
	var sql string = `UPDATE profile_education 
				SET sort_order = $1
				WHERE id = $2 AND user_id = $3`
	tag, err := DbConn.Exec(context.Background(), sql, sortOrder, eduID, userID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reoder profile education: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile education reordering")
	}

	return err
}

func InsertProfileExperience(exp ProfileExperiences) (int, error) {
	var id int
	query := `
		INSERT INTO profile_experiences
		(user_id, experience_type, title, organization, location_text,
		 start_date, end_date, is_current, description, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`

	err := DbConn.QueryRow(
		context.Background(),
		query,
		exp.UserID,
		exp.ExperienceType,
		exp.Title,
		exp.Organization,
		exp.LocationText,
		exp.StartDate,
		exp.EndDate,
		exp.IsCurrent,
		exp.Description,
		exp.SortOrder,
	).Scan(&id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert experience into profile: %v", err)
		return 0, err
	}
	return id, nil
}

func GetProfileExperiences(userID int) ([]ProfileExperiences, error) {
	rows, err := DbConn.Query(context.Background(), `
		SELECT id, user_id, experience_type, title, organization, location_text,
		       start_date, end_date, is_current, description,
		       sort_order, created_at, updated_at
		FROM profile_experiences
		WHERE user_id = $1
		ORDER BY sort_order ASC, start_date DESC
	`, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to grab experience from profile: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []ProfileExperiences
	for rows.Next() {
		var e ProfileExperiences
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.ExperienceType,
			&e.Title,
			&e.Organization,
			&e.LocationText,
			&e.StartDate,
			&e.EndDate,
			&e.IsCurrent,
			&e.Description,
			&e.SortOrder,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to grab experience from profile: %v", err)
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func DeleteProfileExperience(userID, expID int) error {
	tag, err := DbConn.Exec(context.Background(), `
		DELETE FROM profile_experiences
		WHERE id = $1 AND user_id = $2
	`, expID, userID)

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile experience deletion")
	}

	return err
}

func UpdateProfileExperience(exp ProfileExperiences) error {
	query := `
		UPDATE profile_experiences
		SET experience_type = $1,
		    title = $2,
		    organization = $3,
		    location_text = $4,
		    start_date = $5,
		    end_date = $6,
		    is_current = $7,
		    description = $8,
		    sort_order = $9,
		    updated_at = NOW()
		WHERE id = $10 AND user_id = $11
	`

	tag, err := DbConn.Exec(
		context.Background(),
		query,
		exp.ExperienceType,
		exp.Title,
		exp.Organization,
		exp.LocationText,
		exp.StartDate,
		exp.EndDate,
		exp.IsCurrent,
		exp.Description,
		exp.SortOrder,
		exp.ID,
		exp.UserID,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update experience: %v\n", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile experience update")
	}

	return nil
}

func ReorderProfileExperience(userID int, expID int, sortOrder int) error {
	var sql string = `UPDATE profile_experiences
				SET sort_order = $1
				WHERE id = $2 AND user_id = $3`
	tag, err := DbConn.Exec(context.Background(), sql, sortOrder, expID, userID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reoder profile skill: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile skill reordering")
	}

	return err
}

func InsertProfileSkill(skill ProfileSkills) (int, error) {
	var id int
	query := `
		INSERT INTO profile_skills (user_id, skill_name, category, proficiency_label, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := DbConn.QueryRow(
		context.Background(),
		query,
		skill.UserID,
		skill.SkillName,
		skill.Category,
		skill.ProficiencyLabel,
		skill.SortOrder,
	).Scan(&id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert skill into profile: %v", err)
		return 0, err
	}
	return id, nil
}

func GetProfileSkills(userID int) ([]ProfileSkills, error) {
	rows, err := DbConn.Query(context.Background(), `
		SELECT id, user_id, skill_name, category, proficiency_label, sort_order, created_at, updated_at
		FROM profile_skills
		WHERE user_id = $1
		ORDER BY sort_order ASC, id ASC
	`, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to grab skill from profile: %v", err)
		return nil, err
	}
	defer rows.Close()

	var skills []ProfileSkills
	for rows.Next() {
		var s ProfileSkills
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.SkillName,
			&s.Category,
			&s.ProficiencyLabel,
			&s.SortOrder,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to grab skill from profile: %v", err)
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}

func DeleteProfileSkill(userID, skillID int) error {
	tag, err := DbConn.Exec(context.Background(), `
		DELETE FROM profile_skills
		WHERE id = $1 AND user_id = $2
	`, skillID, userID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete skill from profile: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile skill deletion")
	}

	return err
}

func UpdateProfileSkill(skill ProfileSkills) error {
	query := `
		UPDATE profile_skills
		SET skill_name = $1,
		    category = $2,
		    proficiency_label = $3,
		    sort_order = $4,
		    updated_at = NOW()
		WHERE id = $5 AND user_id = $6
	`

	tag, err := DbConn.Exec(
		context.Background(),
		query,
		skill.SkillName,
		skill.Category,
		skill.ProficiencyLabel,
		skill.SortOrder,
		skill.ID,
		skill.UserID,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update skill: %v\n", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile skill update")
	}

	return nil
}

func ReorderProfileSkill(userID int, skillID int, sortOrder int) error {
	var sql string = `UPDATE profile_skills
				SET sort_order = $1
				WHERE id = $2 AND user_id = $3`
	tag, err := DbConn.Exec(context.Background(), sql, sortOrder, skillID, userID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reoder profile skill: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in profile skill reordering")
	}

	return err
}
