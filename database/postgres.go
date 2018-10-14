package database

import (
	"io/ioutil"
	"tpark_db/logger"

	"github.com/jackc/pgx"
)

type Connection struct {
	conn *pgx.Conn
}

// singleton
var DBConn Connection

func (db *Connection) Connect() error {
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
		logger.Logger.Error("err")
		return err
	}
	return nil
}

func (db *Connection) Disconnect() {
	logger.LoggerInfo("Disconnecting database")
	defer db.conn.Close()
	logger.LoggerInfo("Database has been disconnected")
}

func (db *Connection) CreateDB(path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err := db.conn.Exec(string(schema)); err != nil {
		return err
	}

	logger.LoggerInfo("Successfully created tables")
	return nil
}
