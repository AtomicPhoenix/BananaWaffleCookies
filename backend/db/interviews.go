package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Interviews struct {
	ID          int       `json:"id"`
	JobID       int       `json:"job_id"`
	RoundType   string    `json:"round_type"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CreateInterview(interview Interviews) (int, error) {
	var id int

	sqlQuery := `
		INSERT INTO interviews (
			job_id,
			round_type,
			scheduled_at,
			completed_at,
			notes
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := DbConn.QueryRow(
		context.Background(),
		sqlQuery,
		interview.JobID,
		interview.RoundType,
		interview.ScheduledAt,
		interview.CompletedAt,
		interview.Notes,
	).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("Failed to create interview for job_id=%d round_type=%s: %w", interview.JobID, interview.RoundType, err)
	}

	// Activity tracking
	_, err = InsertJobActivity(JobActivity{
		JobID:        int(interview.JobID),
		ActivityType: ActivityInterviewScheduled,
		Description:  fmt.Sprintf("Interview scheduled (%s)", interview.RoundType),
	})
	if err != nil {
		return id, fmt.Errorf("Failed to insert interview activity for job_id=%d interview_id=%d: %w", interview.JobID, id, err)
	}

	return id, nil
}

func DeleteInterview(interviewID int, jobID int, userID int) error {
	// Ensure user owns the job before deleting interview
	isOwner, err := IsJobOwner(jobID, userID)
	if err != nil {
		return fmt.Errorf("Failed to verify job ownership for job_id=%d user_id=%d: %w", jobID, userID, err)
	}
	if !isOwner {
		return fmt.Errorf("Unauthorized delete attempt for interview_id=%d job_id=%d user_id=%d", interviewID, jobID, userID)
	}

	result, err := DbConn.Exec(
		context.Background(),
		`DELETE FROM interviews WHERE id = $1 AND job_id = $2`,
		interviewID,
		jobID,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete interview_id=%d job_id=%d: %w", interviewID, jobID, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected in interview deletion for interview_id=%d job_id=%d", interviewID, jobID)
	}

	// Activity tracking
	_, err = InsertJobActivity(JobActivity{
		JobID:        int(jobID),
		ActivityType: ActivityUpdated,
		Description:  "Interview deleted",
	})
	if err != nil {
		return fmt.Errorf("Failed to insert interview deletion activity for job_id=%d interview_id=%d: %w", jobID, interviewID, err)
	}

	return nil
}

func GetInterviews(jobID int) ([]Interviews, error) {
	query := `
		SELECT 
			id,
			job_id,
			round_type,
			scheduled_at,
			completed_at,
			notes,
			created_at,
			updated_at
		FROM interviews
		WHERE job_id = $1
		ORDER BY scheduled_at ASC;
	`

	rows, err := DbConn.Query(context.Background(), query, jobID)
	if err != nil {
		return nil, fmt.Errorf("Interview query failed for job_id=%d: %w", jobID, err)
	}
	defer rows.Close()

	var interviews []Interviews

	for rows.Next() {
		var i Interviews
		var completedAt sql.NullTime
		var notes sql.NullString

		err := rows.Scan(
			&i.ID,
			&i.JobID,
			&i.RoundType,
			&i.ScheduledAt,
			&completedAt,
			&notes,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Interview scan failed for job_id=%d: %w", jobID, err)
		}

		// Handle nullable fields
		if completedAt.Valid {
			i.CompletedAt = completedAt.Time
		}
		if notes.Valid {
			i.Notes = notes.String
		}

		interviews = append(interviews, i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Interview row iteration failed for job_id=%d: %w", jobID, err)
	}

	return interviews, nil
}
