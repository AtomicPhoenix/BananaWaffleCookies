package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
	Description  string    `json:"description"`
	IsArchived   bool      `json:"is_archived"`
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
	sql_query := `INSERT INTO jobs (user_id, company_name, title, location_text, salary, status, deadline_date, description, is_archived) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
				RETURNING id`
	err := DbConn.QueryRow(context.Background(), sql_query, job.UserID, job.CompanyName, job.Title, job.LocationText, job.Salary, job.Status, job.DeadlineDate, job.Description, false).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert job into database: %v\n", err)
		return -1, err
	}
	return id, err
}

func GetJobs(user_id int, searchQuery string) ([]Job, error) {
	sqlQuery := `
		SELECT id, user_id, company_name, title, location_text, salary, status, deadline_date, description, is_archived, created_at, updated_at FROM jobs WHERE user_id = $1 `

	var (
		rows pgx.Rows
		err  error
	)

	if searchQuery != "" {
		sqlQuery += `
			WHERE company_name ILIKE $2 
			OR title ILIKE $2 
			OR description ILIKE $2
			ORDER BY created_at DESC;`
		searchTerm := "%" + searchQuery + "%"
		rows, err = DbConn.Query(context.Background(), sqlQuery, user_id, searchTerm)
	} else {
		sqlQuery += `ORDER BY created_at DESC;`
		rows, err = DbConn.Query(context.Background(), sqlQuery, user_id)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get jobs from database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job

		var locationText, description sql.NullString
		var deadlineDate sql.NullTime
		var salary sql.NullInt64

		err := rows.Scan(&j.ID, &j.UserID, &j.CompanyName, &j.Title, &locationText, &salary, &j.Status, &deadlineDate, &description, &j.IsArchived, &j.CreatedAt, &j.UpdatedAt)

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
			j.DeadlineDate = deadlineDate.Time
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

func GetJob(job_id int, user_id int) (Job, error) {
	sql_query := `SELECT id, user_id, company_name, title, location_text, salary, status, deadline_date, description, is_archived, created_at, updated_at FROM jobs WHERE id = $1 AND user_id = $2;`

	var job Job
	err := DbConn.QueryRow(context.Background(), sql_query, job_id, user_id).Scan(&job.ID, &job.UserID, &job.CompanyName, &job.Title, &job.LocationText, &job.Salary, &job.Status, &job.DeadlineDate, &job.Description, &job.IsArchived, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get job with id %d from database: %v\n", job_id, err)
		return Job{}, err
	}
	return job, nil
}

func UpdateJob(job Job) error {
	sql_query := `UPDATE jobs 
				SET company_name = $1, title = $2, location_text = $3, salary = $4, status = $5, deadline_date = $6, description = $7
				WHERE id = $8 AND user_id = $9`
	result, err := DbConn.Exec(context.Background(), sql_query, job.CompanyName, job.Title, job.LocationText, job.Salary, job.Status, job.DeadlineDate, job.Description, job.ID, job.UserID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update job: %v\n", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected in update operation (job doesn't exist or is not owned by user)")
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
		fmt.Fprintf(os.Stderr, "Failed to arhive job: %v\n", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected in archive operation (job doesn't exist or is not owned by user)")
	}

	return err
}

func DeleteJob(job Job) error {
	result, err := DbConn.Exec(context.Background(), "DELETE FROM jobs WHERE id=$1 AND user_id = $2", job.ID, job.UserID)
	if err != nil {
		return fmt.Errorf("Failed to delete test job: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected in delete operation (job doesn't exist or is not owned by user)")
	}
	return nil
}
