package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
		fmt.Fprintf(os.Stderr, "Failed to insert interview: %v\n", err)
		return -1, err
	}

	// Activity tracking
	_, err = InsertJobActivity(JobActivity{
		JobID:        int(interview.JobID),
		ActivityType: ActivityInterviewScheduled,
		Description:  fmt.Sprintf("Interview scheduled (%s)", interview.RoundType),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert interview activity: %v\n", err)
		return id, err
	}

	return id, nil
}

func DeleteInterview(interviewID int, jobID int, userID int) error {
	// Ensure user owns the job before deleting interview
	isOwner, err := IsJobOwner(jobID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("user does not own this job")
	}

	result, err := DbConn.Exec(
		context.Background(),
		`DELETE FROM interviews WHERE id = $1 AND job_id = $2`,
		interviewID,
		jobID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete interview: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected (interview not found)")
	}

	// Activity tracking
	_, err = InsertJobActivity(JobActivity{
		JobID:        int(jobID),
		ActivityType: ActivityUpdated,
		Description:  "Interview deleted",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert interview delete activity: %v\n", err)
		return err
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
		fmt.Fprintf(os.Stderr, "Failed to query interviews: %v\n", err)
		return nil, err
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
			fmt.Fprintf(os.Stderr, "Failed to scan interview: %v\n", err)
			return nil, err
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
		fmt.Fprintf(os.Stderr, "Row iteration error: %v\n", err)
		return nil, err
	}

	return interviews, nil
}
