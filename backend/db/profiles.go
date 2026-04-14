package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Profile struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	FirstName          string    `json:"first_name,omitempty"`
	LastName           string    `json:"last_name,omitempty"`
	Phone              string    `json:"phone,omitempty"`
	Location           string    `json:"location,omitempty"`
	City               string    `json:"city,omitempty"`
	State              string    `json:"state,omitempty"`
	Country            string    `json:"country,omitempty"`
	Headline           string    `json:"headline,omitempty"`
	LinkedinURL        string    `json:"linkedin_url,omitempty"`
	PortfolioURL       string    `json:"portfolio_url,omitempty"`
	Summary            string    `json:"summary,omitempty"`
	PreferredCity      string    `json:"preferred_city,omitempty"`
	PreferredState     string    `json:"preferred_state,omitempty"`
	PreferredRole      string    `json:"preferred_role,omitempty"`
	PreferredSalaryMin int       `json:"preferred_salary_min,omitempty"`
	PreferredSalaryMax int       `json:"preferred_salary_max,omitempty"`
	WorkMode           string    `json:"work_mode,omitempty"`
	CompletionPercent  int       `json:"completion_percent"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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
	query := `UPDATE profiles SET
		first_name = $1,
		last_name = $2,
		phone = $3,
		location = $4,
		city = $5,
		state = $6,
		country = $7,
		headline = $8,
		linkedin_url = $9,
		portfolio_url = $10,
		summary = $11,
		preferred_city = $12,
		preferred_state = $13,
		preferred_role = $14,
		preferred_salary_min = $15,
		preferred_salary_max = $16,
		work_mode = $17,
		completion_percent = $18,
		updated_at = NOW()
	WHERE user_id = $19`

	if profile.CompletionPercent == 0 {
		profile.SetCompletionPercent()
	}

	_, err := DbConn.Exec(
		context.Background(),
		query,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.Location,
		profile.City,
		profile.State,
		profile.Country,
		profile.Headline,
		profile.LinkedinURL,
		profile.PortfolioURL,
		profile.Summary,
		profile.PreferredCity,
		profile.PreferredState,
		profile.PreferredRole,
		profile.PreferredSalaryMin,
		profile.PreferredSalaryMax,
		profile.WorkMode,
		profile.CompletionPercent,
		profile.UserID,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update profile: %v\n", err)
		return err
	}
	return nil
}

func GetProfile(uid int) (Profile, error) {
	var profile Profile

	var firstName, lastName, phone, location, city, state, country sql.NullString
	var headline, linkedinURL, portfolioURL, summary sql.NullString
	var preferredCity, preferredState, preferredRole, workMode sql.NullString
	var preferredSalaryMin, preferredSalaryMax sql.NullInt64

	err := DbConn.QueryRow(context.Background(), `
		SELECT
			first_name,
			last_name,
			phone,
			location,
			city,
			state,
			country,
			headline,
			linkedin_url,
			portfolio_url,
			summary,
			preferred_city,
			preferred_state,
			preferred_role,
			preferred_salary_min,
			preferred_salary_max,
			work_mode,
			completion_percent
		FROM profiles
		WHERE user_id = $1
	`, uid).Scan(
		&firstName,
		&lastName,
		&phone,
		&location,
		&city,
		&state,
		&country,
		&headline,
		&linkedinURL,
		&portfolioURL,
		&summary,
		&preferredCity,
		&preferredState,
		&preferredRole,
		&preferredSalaryMin,
		&preferredSalaryMax,
		&workMode,
		&profile.CompletionPercent,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get profile: %v\n", err)
		return Profile{}, err
	}

	profile.UserID = uid
	profile.FirstName = extractValue(firstName)
	profile.LastName = extractValue(lastName)
	profile.Phone = extractValue(phone)
	profile.Location = extractValue(location)
	profile.City = extractValue(city)
	profile.State = extractValue(state)
	profile.Country = extractValue(country)
	profile.Headline = extractValue(headline)
	profile.LinkedinURL = extractValue(linkedinURL)
	profile.PortfolioURL = extractValue(portfolioURL)
	profile.Summary = extractValue(summary)
	profile.PreferredCity = extractValue(preferredCity)
	profile.PreferredState = extractValue(preferredState)
	profile.PreferredRole = extractValue(preferredRole)
	profile.WorkMode = extractValue(workMode)

	if preferredSalaryMin.Valid {
		profile.PreferredSalaryMin = int(preferredSalaryMin.Int64)
	}
	if preferredSalaryMax.Valid {
		profile.PreferredSalaryMax = int(preferredSalaryMax.Int64)
	}

	return profile, nil
}

func extractValue(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}

func (profile *Profile) SetCompletionPercent() {
	profile.CompletionPercent = profile.getProfileCompletionPercent()
}

// Calculate what percent of the profile is filled out (returns an int from 0 to 100)
func (profile *Profile) getProfileCompletionPercent() int {
	var filledFields int = 0
	var numFields int = 17

	if profile.FirstName != "" {
		filledFields++
	}
	if profile.LastName != "" {
		filledFields++
	}
	if profile.Phone != "" {
		filledFields++
	}
	if profile.Location != "" {
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
	if profile.Headline != "" {
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
	if profile.PreferredCity != "" {
		filledFields++
	}
	if profile.PreferredState != "" {
		filledFields++
	}
	if profile.PreferredRole != "" {
		filledFields++
	}
	if profile.PreferredSalaryMin > 0 {
		filledFields++
	}
	if profile.PreferredSalaryMax > 0 {
		filledFields++
	}
	if profile.WorkMode != "" {
		filledFields++
	}
	return int((float32(filledFields) / float32(numFields)) * 100)
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
