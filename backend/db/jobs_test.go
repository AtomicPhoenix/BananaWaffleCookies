package db

import (
	"context"
	"testing"
	"time"
)

// Does not delete test user since job is dependent on it
func createTestJob(t *testing.T) (User, Job) {
	test_user := createTestUser(t)
	test_job := Job{
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
	return test_user, test_job
}

func TestCreateJob(t *testing.T) {
	test_user, _ := createTestJob(t)
	deleteTestUser(t, test_user.Id)
}

func TestRetrieveJob(t *testing.T) {
	test_user, test_job := createTestJob(t)
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	retrieved_job, err := GetJob(test_job.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve job: %v", err)
	}

	if retrieved_job.CompanyName != test_job.CompanyName {
		t.Errorf("Failed to retrieve job with correct information: expected %s, got %s", test_job.CompanyName, retrieved_job.CompanyName)
	}
}

func TestUpdateJob(t *testing.T) {
	test_user, test_job := createTestJob(t)
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

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
	test_user, test_job := createTestJob(t)
	t.Cleanup(func() {
		deleteTestUser(t, test_user.Id)
	})

	_, err := DbConn.Exec(context.Background(), "DELETE FROM jobs WHERE id=$1", test_job.ID)
	if err != nil {
		t.Fatalf("Failed to delete test job: %v", err)
	}
}
