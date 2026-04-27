package db

import (
	"context"
	"fmt"
	"time"
)

type FollowUpTask struct {
	ID          int        `json:"id"`
	JobID       int        `json:"job_id"`
	Title       string     `json:"title"`
	Notes       string     `json:"notes,omitempty"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	RemindAt    *time.Time `json:"remind_at,omitempty"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CREATE
func CreateFollowUp(jobID int, f FollowUpTask) (FollowUpTask, error) {
	query := `
		INSERT INTO follow_up_tasks (job_id, title, notes, due_at, remind_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, job_id, title, notes, due_at, remind_at,
		          is_completed, completed_at, created_at, updated_at;
	`

	var created FollowUpTask

	err := DbConn.QueryRow(context.Background(), query,
		jobID,
		f.Title,
		f.Notes,
		f.DueAt,
		f.RemindAt,
	).Scan(
		&created.ID,
		&created.JobID,
		&created.Title,
		&created.Notes,
		&created.DueAt,
		&created.RemindAt,
		&created.IsCompleted,
		&created.CompletedAt,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return FollowUpTask{}, fmt.Errorf("Failed to create follow_up_task for job_id=%d title=%s: %w", jobID, f.Title, err)
	}

	return created, nil
}

// GET ALL
func GetFollowUps(jobID int) ([]FollowUpTask, error) {
	query := `
		SELECT id, job_id, title, notes, due_at, remind_at,
		       is_completed, completed_at, created_at, updated_at
		FROM follow_up_tasks
		WHERE job_id = $1
		ORDER BY due_at ASC;
	`

	rows, err := DbConn.Query(context.Background(), query, jobID)
	if err != nil {
		return nil, fmt.Errorf("FollowUpTask query failed for job_id=%d: %w", jobID, err)
	}
	defer rows.Close()

	var tasks []FollowUpTask

	for rows.Next() {
		var t FollowUpTask

		err := rows.Scan(
			&t.ID,
			&t.JobID,
			&t.Title,
			&t.Notes,
			&t.DueAt,
			&t.RemindAt,
			&t.IsCompleted,
			&t.CompletedAt,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("FollowUpTask scan failed for job_id=%d: %w", jobID, err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

// UPDATE
func UpdateFollowUp(id int, jobID int, f FollowUpTask) (FollowUpTask, error) {
	query := `
		UPDATE follow_up_tasks
		SET title = $1,
		    notes = $2,
		    due_at = $3,
		    remind_at = $4,
		    is_completed = $5,
		    completed_at = CASE
		        WHEN $5 = true THEN NOW()
		        ELSE NULL
		    END,
		    updated_at = NOW()
		WHERE id = $6 AND job_id = $7
		RETURNING id, job_id, title, notes, due_at, remind_at,
		          is_completed, completed_at, created_at, updated_at;
	`

	var updated FollowUpTask

	err := DbConn.QueryRow(context.Background(), query,
		f.Title,
		f.Notes,
		f.DueAt,
		f.RemindAt,
		f.IsCompleted,
		id,
		jobID,
	).Scan(
		&updated.ID,
		&updated.JobID,
		&updated.Title,
		&updated.Notes,
		&updated.DueAt,
		&updated.RemindAt,
		&updated.IsCompleted,
		&updated.CompletedAt,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return FollowUpTask{}, fmt.Errorf("Failed to update follow_up_task id=%d job_id=%d: %w", id, jobID, err)
	}

	return updated, nil
}

// DELETE
func DeleteFollowUp(id int, jobID int) error {
	query := `
		DELETE FROM follow_up_tasks
		WHERE id = $1 AND job_id = $2;
	`

	_, err := DbConn.Exec(context.Background(), query, id, jobID)
	if err != nil {
		return fmt.Errorf("Failed to delete follow_up_task id=%d job_id=%d: %w", id, jobID, err)
	}

	return nil
}
