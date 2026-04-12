package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DbConn *pgxpool.Pool

func get_db_connection_string() string {
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_addr := os.Getenv("DB_ADDR")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")

	if db_name == "" {
		panic("DB_NAME is not set")
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_addr, db_port, db_name)
}

func InitDB() error {
	fmt.Println("Initializing Database")
	var err error
	DbConn, err = pgxpool.New(context.Background(), get_db_connection_string())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	fmt.Println("Successfully connected to database.")
	return nil
}

func main() {
	err := InitDB()
	if err != nil {
		os.Exit(1)
	}
	defer DbConn.Close()
}
