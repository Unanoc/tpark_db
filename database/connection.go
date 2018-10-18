package database

import (
	"io/ioutil"
	"log"

	"github.com/jackc/pgx"
)

var db *pgx.ConnPool

var pgxConfig = pgx.ConnConfig{
	User:              "forum",
	Password:          "forum",
	Host:              "localhost",
	Port:              5432,
	Database:          "forum",
	TLSConfig:         nil,
	UseFallbackTLS:    false,
	FallbackTLSConfig: nil,
}

const dataBaseSchema = "./sql/create_tables.sql"
const removeDataBase = "./sql/drop_tables.sql"

func Connect() {
	var err error
	if db, err = pgx.NewConnPool( // creates a new ConnPool. config.ConnConfig is passed through to Connect directly.
		pgx.ConnPoolConfig{
			ConnConfig:     pgxConfig,
			MaxConnections: 8,
		}); err != nil {
		log.Fatalln(err) // Fatalln is equivalent to Println() followed by a call to os.Exit(1)
	}

	// debug
	if err = dropTablesIfExist(); err != nil {
		log.Println(err)
	}
	log.Println("SQL Schema was dropped successfully")
	//

	if err = createTables(); err != nil {
		log.Println(err)
	}
	log.Println("SQL Schema was initialized successfully")
}

func Disconnect() {
	db.Close()
}

func createTables() error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	content, err := ioutil.ReadFile(dataBaseSchema)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := tx.Exec(string(content)); err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func dropTablesIfExist() error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	content, err := ioutil.ReadFile(removeDataBase)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := tx.Exec(string(content)); err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func StartTransaction() *pgx.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	return tx
}

func CommitTransaction(tx *pgx.Tx) {
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err)
	}
}
