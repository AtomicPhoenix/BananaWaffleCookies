package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Job struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	CompanyName  string    `json:"company_name"`
	Title        string    `json:"title"`
	LocationText string    `json:"location_text"`
	Salary       int       `json:"salary"`
	Status       string    `json:"status"`
	DeadlineDate string    `json:"deadline_date"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Job_Activity struct {
	job_id        int
	activity_type string // in ( 'created', 'updated', 'status_changed', 'applied', 'note_added', 'document_linked', 'document_unlinked')
	time          time.Time
	description   string
	metadata      string
}

func CreateJob(job Job) (int, error) {
	var id int
	sql_query := `INSERT INTO jobs (user_id, company_name, title, location_text, salary, status, deadline_date, description) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
				RETURNING id`
	err := DbConn.QueryRow(context.Background(), sql_query, job.UserID, job.CompanyName, job.Title, job.LocationText, job.Salary, job.Status, job.DeadlineDate, job.Description).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert job into database: %v\n", err)
		return -1, err
	}
	return id, err
}

func GetAllJobs() ([]Job, error) {
	sql_query := `SELECT id, user_id, company_name, title, location_text, salary, status, deadline_date, description, created_at, updated_at FROM jobs ORDER BY created_at DESC;`
	rows, err := DbConn.Query(context.Background(), sql_query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get jobs from database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job

		var locationText, deadlineDate, description sql.NullString
		var salary sql.NullInt64

		err := rows.Scan(&j.ID, &j.UserID, &j.CompanyName, &j.Title, &locationText, &salary, &j.Status, &deadlineDate, &description, &j.CreatedAt, &j.UpdatedAt)

		if err != nil {
			return nil, err
		}

		// Convert nullable fields
		if locationText.Valid {
			j.LocationText = locationText.String
		}
		if salary.Valid {
			j.Salary = int(salary.Int64)
		}
		if deadlineDate.Valid {
			j.DeadlineDate = deadlineDate.String
		}
		if description.Valid {
			j.Description = description.String
		}

		jobs = append(jobs, j)
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get jobs from database: %v\n", err)
		return nil, err
	}
	return jobs, nil
}

func UpdateJob(job Job) error {
	sql_query := `UPDATE jobs 
				SET company_name = $1, title = $2, location_text = $3, salary = $4, status = $5, deadline_date = $6, description = $7
				WHERE id = $8`
	_, err := DbConn.Exec(context.Background(), sql_query, job.CompanyName, job.Title, job.LocationText, job.Salary, job.Status, job.DeadlineDate, job.Description, job.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update job: %v\n", err)
		return err
	}
	return err
}
