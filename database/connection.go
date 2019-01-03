package database

import "github.com/jackc/pgx"

// DB is a main instance connection.
type DataBase struct {
	Conn       *pgx.ConnPool
	SchemaPath string
}

// DB is the global instance of DataBase structure.
var DB DataBase

// Connect creates a connection with db.
func (db *DataBase) Connect(psqlURI string) error {
	pgxConfig, err := pgx.ParseURI(psqlURI)
	if err != nil {
		return err
	}

	if db.Conn, err = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgxConfig,
		}); err != nil {
		return err
	}

	return nil
}

// Disconnect closes a connection.
func (db *DataBase) Disconnect() {
	db.Conn.Close()
}
