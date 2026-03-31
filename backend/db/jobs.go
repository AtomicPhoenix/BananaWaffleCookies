package db

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Job struct {
	user_id       int
	company_name  string
	title         string
	location_text string
	posting_url   string
	status        string // in ('interested', 'applied', 'interview', 'offer', 'rejected', 'archived')
	deadline_date time.Time
	notes         string
}

type Job_Activity struct {
	job_id        int
	activity_type string // in ( 'created', 'updated', 'status_changed', 'applied', 'note_added', 'document_linked', 'document_unlinked')
	time          time.Time
	description   string
	metadata      string
}

func GetJob(title string) int {
	var id int
	err := DbConn.QueryRow(context.Background(), "SELECT id FROM jobs WHERE title=$1", title).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get job: %v\n", err)
	}
	return id
}

func CreateJob(job Job) (bool, int) {
	var id int
	sql := `INSERT INTO jobs (user_id, company_name, title, location_text, posting_url, status, deadline_date, notes) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
				RETURNING id`
	err := DbConn.QueryRow(context.Background(), sql, job.user_id, job.company_name, job.title, job.location_text, job.posting_url, job.status, job.deadline_date, job.notes).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert job into database: %v\n", err)
		return false, -1
	}
	return true, id
}
