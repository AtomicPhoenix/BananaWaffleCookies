package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Job struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	CompanyName  string    `json:"company_name"`
	Title        string    `json:"title"`
	LocationText string    `json:"location_text"`
	Salary       int       `json:"salary"`
	Status       string    `json:"status"`
	DeadlineDate time.Time `json:"deadline_date"`
	Notes        string    `json:"notes"`
	Description  string    `json:"description"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Struct for tracking updates to a job
type JobActivity struct {
	ID           int          `json:"id"`
	JobID        int          `json:"job_id"`
	ActivityType ActivityType `json:"activity_type"`
	ActivityAt   time.Time    `json:"activity_at"` // When the activity occured
	Description  string       `json:"description,omitempty"`
}

// Valid activity types
type ActivityType string

const (
	ActivityCreated            ActivityType = "created"
	ActivityUpdated            ActivityType = "updated"
	ActivityStatusChanged      ActivityType = "status_changed"
	ActivityApplied            ActivityType = "applied"
	ActivityNoteAdded          ActivityType = "note_added"
	ActivityDocumentLinked     ActivityType = "document_linked"
	ActivityDocumentUnlinked   ActivityType = "document_unlinked"
	ActivityInterviewScheduled ActivityType = "interview_scheduled"
	ActivityInterviewCompleted ActivityType = "interview_completed"
	ActivityFollowUpCreated    ActivityType = "follow_up_created"
	ActivityFollowUpCompleted  ActivityType = "follow_up_completed"
	ActivityOutcome            ActivityType = "outcome"
)

func CreateJob(job Job) (int, error) {
	var id int
	sql_query := `INSERT INTO jobs (user_id, company_name, title, location_text, salary, status, deadline_date, notes, description, is_archived) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
				RETURNING id`
	err := DbConn.QueryRow(context.Background(), sql_query, job.UserID, job.CompanyName, job.Title, job.LocationText, job.Salary, job.Status, job.DeadlineDate, job.Notes, job.Description, false).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("Failed to insert job for user_id=%d: %w", job.UserID, err)

	}

	_, err = InsertJobActivity(JobActivity{
		JobID:        id,
		ActivityType: ActivityCreated,
		Description:  "Job created",
	})
	if err != nil {
		return id, fmt.Errorf("Failed to insert job creation activity for job_id=%d: %w", id, err)
	}

	return id, err
}

func GetJobs(user_id int, searchQuery string) ([]Job, error) {
	sqlQuery := `
		SELECT id, user_id, company_name, title, location_text, salary, status, deadline_date, notes, description, is_archived, created_at, updated_at FROM jobs WHERE user_id = $1 `

	var (
		rows pgx.Rows
		err  error
	)

	if searchQuery != "" {
		sqlQuery += `
			AND (company_name ILIKE $2 
			OR title ILIKE $2 
			OR description ILIKE $2)
			ORDER BY created_at DESC;`
		searchTerm := "%" + searchQuery + "%"
		rows, err = DbConn.Query(context.Background(), sqlQuery, user_id, searchTerm)
	} else {
		sqlQuery += `ORDER BY created_at DESC;`
		rows, err = DbConn.Query(context.Background(), sqlQuery, user_id)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to query jobs for user_id=%d: %w", user_id, err)
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job

		var locationText, description, notes sql.NullString
		var deadlineDate sql.NullTime
		var salary sql.NullInt64

		err := rows.Scan(&j.ID, &j.UserID, &j.CompanyName, &j.Title, &locationText, &salary, &j.Status, &deadlineDate, &notes, &description, &j.IsArchived, &j.CreatedAt, &j.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan jobs for user_id=%d: %w", user_id, err)
		}

		// Convert nullable fields
		if locationText.Valid {
			j.LocationText = locationText.String
		}
		if salary.Valid {
			j.Salary = int(salary.Int64)
		}
		if deadlineDate.Valid {
			j.DeadlineDate = deadlineDate.Time
		}
		if notes.Valid {
			j.Notes = notes.String
		}
		if description.Valid {
			j.Description = description.String
		}

		jobs = append(jobs, j)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to iterate jobs result for user_id=%d: %w", user_id, err)
	}
	return jobs, nil
}

func GetJob(job_id int, user_id int) (Job, error) {
	sql_query := `SELECT id, user_id, company_name, title, location_text, salary, status, deadline_date, notes, description, is_archived, created_at, updated_at FROM jobs WHERE id = $1 AND user_id = $2;`

	var job Job
	err := DbConn.QueryRow(context.Background(), sql_query, job_id, user_id).Scan(&job.ID, &job.UserID, &job.CompanyName, &job.Title, &job.LocationText, &job.Salary, &job.Status, &job.DeadlineDate, &job.Notes, &job.Description, &job.IsArchived, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return Job{}, fmt.Errorf("Failed to get job job_id=%d user_id=%d: %w", job_id, user_id, err)
	}
	return job, nil
}

func UpdateJob(job Job) error {
	oldJob, err := GetJob(job.ID, job.UserID)
	if err != nil {
		return fmt.Errorf("Failed to update job job_id=%d user_id=%d: %w. Failed to get old job information", job.ID, job.UserID, err)
	}

	sql_query := `UPDATE jobs 
				SET company_name = $1, title = $2, location_text = $3, salary = $4, status = $5, deadline_date = $6, notes = $7, description = $8
				WHERE id = $9 AND user_id = $10`
	result, err := DbConn.Exec(context.Background(), sql_query, job.CompanyName, job.Title, job.LocationText, job.Salary, job.Status, job.DeadlineDate, job.Notes, job.Description, job.ID, job.UserID)
	if err != nil {
		return fmt.Errorf("Failed to update job job_id=%d user_id=%d: %w", job.ID, job.UserID, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected when updating job job_id=%d user_id=%d", job.ID, job.UserID)
	}

	jobActivity, err := createJobActivity(oldJob, job)
	if err != nil {
		return fmt.Errorf("Failed to insert job update activity for job_id=%d: %w", job.ID, err)
	}

	_, err = InsertJobActivity(jobActivity)
	if err != nil {
		return fmt.Errorf("Failed to insert job update activity for job_id=%d: %w", job.ID, err)
	}

	return err
}

func ArchiveJob(job Job) error {
	return setArchive(job, true)
}

func UnarchiveJob(job Job) error {
	return setArchive(job, false)
}

func setArchive(job Job, is_archived bool) error {
	sql_query := `UPDATE jobs 
				SET is_archived = $1
				WHERE id = $2 AND user_id = $3`
	result, err := DbConn.Exec(context.Background(), sql_query, is_archived, job.ID, job.UserID)
	if err != nil {
		return fmt.Errorf("Failed to update job archive state job_id=%d user_id=%d: %w", job.ID, job.UserID, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected when updating archive state job_id=%d user_id=%d", job.ID, job.UserID)
	}

	var desc string
	if is_archived {
		desc = "Job archived"
	} else {
		desc = "Job unarchived"
	}

	_, err = InsertJobActivity(JobActivity{
		JobID:        job.ID,
		ActivityType: ActivityStatusChanged,
		Description:  desc,
	})

	if err != nil {
		return fmt.Errorf("Failed to insert job activity for archive state change job_id=%d: %w", job.ID, err)
	}

	return nil
}

func DeleteJob(job Job) error {
	result, err := DbConn.Exec(context.Background(), "DELETE FROM jobs WHERE id=$1 AND user_id = $2", job.ID, job.UserID)
	if err != nil {
		return fmt.Errorf("Failed to delete job job_id=%d user_id=%d: %w", job.ID, job.UserID, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected when deleting job job_id=%d user_id=%d", job.ID, job.UserID)
	}
	return nil
}

// JOB ACTIVITY HELPERS
func InsertJobActivity(activity JobActivity) (JobActivity, error) {
	query := `
		INSERT INTO job_activities (
			job_id,
			activity_type,
			description
		)
		VALUES ($1, $2, $3)
		RETURNING id, activity_at
	`

	err := DbConn.QueryRow(
		context.Background(),
		query,
		activity.JobID,
		activity.ActivityType,
		activity.Description,
	).Scan(&activity.ID, &activity.ActivityAt)

	if err != nil {
		return JobActivity{}, fmt.Errorf("Failed to insert job activity for job_id=%d: %w", activity.JobID, err)
	}

	return activity, nil
}

func createJobActivity(oldJob, newJob Job) (JobActivity, error) {
	if oldJob.ID != newJob.ID {
		return JobActivity{}, fmt.Errorf("Failed to create job activity: oldJob.ID != newJob.ID")
	}

	activity := JobActivity{
		JobID:        newJob.ID,
		ActivityType: ActivityUpdated,
		Description:  "Job Updated:",
	}

	var changes []string

	if oldJob.CompanyName != newJob.CompanyName {
		changes = append(changes, fmt.Sprintf("Company name updated to %s", newJob.CompanyName))
	}

	if oldJob.Title != newJob.Title {
		changes = append(changes, fmt.Sprintf("Title updated to %s", newJob.Title))
	}

	if oldJob.LocationText != newJob.LocationText {
		changes = append(changes, fmt.Sprintf("Location updated to %s", newJob.LocationText))
	}

	if oldJob.Salary != newJob.Salary {
		changes = append(changes, fmt.Sprintf("Salary updated to %d", newJob.Salary))
	}

	if oldJob.Status != newJob.Status {
		changes = append(changes, fmt.Sprintf("Status updated to %s", newJob.Status))
	}

	if !oldJob.DeadlineDate.Equal(newJob.DeadlineDate) {
		changes = append(changes, fmt.Sprintf("Deadline updated to %s", newJob.DeadlineDate.Format(time.RFC3339)))
	}

	if oldJob.Description != newJob.Description {
		changes = append(changes, "Description updated")
	}

	if oldJob.IsArchived != newJob.IsArchived {
		changes = append(changes, fmt.Sprintf("Archived status changed to %t", newJob.IsArchived))
	}

	if len(changes) == 0 {
		activity.Description = "Job Updated: no changes detected"
		return activity, nil
	}

	activity.Description = activity.Description + "\n\t- " + strings.Join(changes, "\n\t- ")

	return activity, nil
}

func GetJobActivities(jobID int) ([]JobActivity, error) {
	query := `
		SELECT
			id,
			job_id,
			activity_type,
			activity_at,
			description
		FROM job_activities
		WHERE job_id = $1
		ORDER BY activity_at DESC
	`

	rows, err := DbConn.Query(context.Background(), query, jobID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query job activities for job_id=%d: %w", jobID, err)
	}
	defer rows.Close()

	var activities []JobActivity

	for rows.Next() {
		var activity JobActivity
		var description sql.NullString

		err := rows.Scan(
			&activity.ID,
			&activity.JobID,
			&activity.ActivityType,
			&activity.ActivityAt,
			&description,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan job activity for job_id=%d: %w", jobID, err)
		}

		if description.Valid {
			activity.Description = description.String
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func IsJobOwner(jobID int, userID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM jobs
			WHERE id = $1 AND user_id = $2
		);
	`

	var exists bool

	err := DbConn.QueryRow(
		context.Background(),
		query,
		jobID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("Failed to verify job ownership job_id=%d user_id=%d: %w", jobID, userID, err)
	}

	return exists, nil
}

func UpdateJobCompanyNotes(jobID int, userID int, notes string) error {
	query := `
		UPDATE jobs
		SET notes = $1
		WHERE id = $2 AND user_id = $3
	`

	res, err := DbConn.Exec(context.Background(), query, notes, jobID, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}
