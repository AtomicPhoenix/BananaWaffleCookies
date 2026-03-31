package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DbConn *pgx.Conn

func get_db_connection_string() string {
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_addr := os.Getenv("DB_ADDR")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_addr, db_port, db_name)
}

func init() {
	fmt.Printf("-----------------------\nInitializing Database\n")
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load env file: %v\n", err)
		os.Exit(1)
	}

	DbConn, err = pgx.Connect(context.Background(), get_db_connection_string())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err = DbConn.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to database.")

	checkTables()

	fmt.Printf("Database setup completed!\n-----------------------\n")
}

func main() {
	defer func() {
		if err := DbConn.Close(context.Background()); err != nil {
			log.Fatalf("Error in closing database: %s", err)
		}
	}()
}

func checkTables() {
	var result bool
	tables := []string{"users", "profiles", "jobs", "job_activities", "documents", "document_versions"}
	for _, table := range tables {
		err := DbConn.QueryRow(context.Background(), "SELECT EXISTS ( SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = '$1');", table).Scan(&result)
		if err != nil || result == false {
			fmt.Fprintf(os.Stderr, "Table does not exist: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("All tables exist!")
}
