package db

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	os.Exit(m.Run())
}

func TestInitDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Errorf(`Failed to init database: %v`, err)
	}
}

func TestPingDB(t *testing.T) {
	err := DbConn.Ping(context.Background())
	if err != nil {
		t.Errorf(`Failed to ping database: %v`, err)
	}
}

func TestTableExistence(t *testing.T) {
	var result bool
	tables := []string{"users", "profiles", "jobs", "job_activities", "documents", "document_links"}
	for _, table := range tables {
		err := DbConn.QueryRow(context.Background(), "SELECT EXISTS ( SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = $1);", table).Scan(&result)
		if err != nil || result == false {
			t.Errorf("Table %s does not exist: %v\n", table, err)
		}
	}
}
