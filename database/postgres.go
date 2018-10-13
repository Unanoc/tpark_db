package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

type Database struct {
	conn *pgx.Conn
}

func (db *Database) Connect() error {
	runtimeParams := make(map[string]string)
	runtimeParams["application_name"] = "tpark_db"
	connConfig := pgx.ConnConfig{
		User:              "forum",
		Password:          "forum",
		Host:              "localhost",
		Port:              5432,
		Database:          "forum",
		TLSConfig:         nil,
		UseFallbackTLS:    false,
		FallbackTLSConfig: nil,
		RuntimeParams:     runtimeParams,
	}
	conn, err := pgx.Connect(connConfig)
	db.conn = conn
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		return err
	}
	return nil
}

func (db *Database) Disconnect() {
	fmt.Println("Disconnecting database")
	defer db.conn.Close()
	fmt.Println("Database has been disconnected")
}

func (db Database) CreateDB() error {

	_, err := db.conn.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		"nickname" CITEXT UNIQUE PRIMARY KEY,
		"email"    CITEXT UNIQUE NOT NULL,
		"fullname" CITEXT NOT NULL,
		"about"    TEXT
	  );
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create users table: %v\n", err)
		return err
	}
	fmt.Printf("Successfully created users table\n")
	return nil
}
