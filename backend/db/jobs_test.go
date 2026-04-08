package db

import (
	"context"
	"testing"
	"time"
)

var test_user User
var test_job Job

func TestCreateJob(t *testing.T) {
	test_user = createTestUser(t)
	test_job = Job{
		UserID:       test_user.Id,
		CompanyName:  "Aperture Labs",
		Title:        "Tester",
		Status:       "applied",
		DeadlineDate: time.Now(),
		Description:  "Portal Gun tester",
	}

	jobID, err := CreateJob(test_job)
	test_job.ID = jobID
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}
}

func TestRetrieveJob(t *testing.T) {
	retrieved_job, err := GetJob(test_job.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve job: %v", err)
	}

	if retrieved_job.CompanyName != test_job.CompanyName {
		t.Errorf("Failed to retrieve job with correct information: expected %s, got %s", test_job.CompanyName, retrieved_job.CompanyName)
	}
}

func TestUpdateJob(t *testing.T) {
	test_job.CompanyName = "Aperture Labs"
	test_job.Title = "Cave Johnson"
	test_job.Status = "applied"
	test_job.DeadlineDate = time.Now()
	test_job.Description = "Owner, CEO, Visionary"

	err := UpdateJob(test_job)
	if err != nil {
		t.Fatalf("Failed to update test user: %v", err)
	}

	retrieved_job, err := GetJob(test_job.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated job: %v", err)
	}

	if retrieved_job.CompanyName != test_job.CompanyName {
		t.Errorf("Failed to retrieve job with correct information: expected %s, got %s", test_job.CompanyName, retrieved_job.CompanyName)
	}
}

func TestDeleteJob(t *testing.T) {
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	_, err := DbConn.Exec(context.Background(), "DELETE FROM jobs WHERE id=$1", test_job.ID)
	if err != nil {
		t.Fatalf("failed to delete test user: %v", err)
	}
}
