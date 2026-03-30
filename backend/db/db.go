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
	err := godotenv.Load() 
	if err != nil {
		log.Fatal(err)
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
}

func main() {
	defer func() {
		if err := DbConn.Close(context.Background()); err != nil {
			log.Fatalf("Error in closing database: %s", err)
		}
	}()
}
