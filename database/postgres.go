// package database

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"github.com/jackc/pgx"
// )

// type Connection struct {
// 	conn *pgx.Conn
// }

// // singleton
// var DBConn Connection

// func (db *Connection) Connect() error {
// 	runtimeParams := make(map[string]string)
// 	runtimeParams["application_name"] = "tpark_db"
// 	connConfig := pgx.ConnConfig{
// 		User:              "forum",
// 		Password:          "forum",
// 		Host:              "localhost",
// 		Port:              5432,
// 		Database:          "forum",
// 		TLSConfig:         nil,
// 		UseFallbackTLS:    false,
// 		FallbackTLSConfig: nil,
// 		RuntimeParams:     runtimeParams,
// 	}
// 	conn, err := pgx.Connect(connConfig)
// 	DBConn.conn = conn

// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
// 		return err
// 	}
// 	return nil
// }

// func (db *Connection) Disconnect() {
// 	fmt.Println("Disconnecting database")
// 	defer db.conn.Close()
// 	fmt.Println("Database has been disconnected")
// }

// func (db *Connection) CreateDB(path string) error {
// 	tx, err := StartTransaction()
// 	defer tx.Rollback()

// 	schema, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := db.conn.Exec(string(schema)); err != nil {
// 		return err
// 	}
// 	CommitTransaction(tx)

// 	fmt.Printf("Successfully created tables\n")
// 	return nil
// }

// func StartTransaction() (*pgx.Tx, error) {
// 	tx, err := DBConn.conn.Begin()
// 	if err != nil {
// 		return tx, err
// 	}
// 	return tx, nil
// }

// func CommitTransaction(tx *pgx.Tx) {
// 	if err := tx.Commit(); err != nil {
// 		tx.Rollback()
// 	}
// }
