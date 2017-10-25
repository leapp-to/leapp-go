// Package db provides operations on SQLite DB.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	// Register sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Database type is sql.DB
type Database struct {
	conn *sql.DB
}

// Record type represents Database record.
type Record struct {
	UID, Name, Host, DataType, Data string
}

// New initializes database, connect to it and return the connection
func New(dbDir, dbName string) (*Database, error) {
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}
	conn, err := sql.Open("sqlite3", filepath.Join(dbDir, dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the databse: %v", err)
	}
	return &Database{conn: conn}, nil
}

// Close closes the connection to the database.
func (db *Database) Close() error {
	return db.conn.Close()
}

// Exec executes given SQL statement.
func (db *Database) Exec(sqlStatement string) error {
	_, err := db.conn.Exec(sqlStatement)
	return err
}

// SelectByUIDNameTypeAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDNameTypeAndHost(uid, name, dataType, host string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND data_type = ? AND host = ?"
	return db.internalSelect(whereClause, uid, name, dataType, host)
}

// SelectByUIDNameAndType prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDNameAndType(uid, name, dataType string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND data_type = ?"
	return db.internalSelect(whereClause, uid, name, dataType)
}

// SelectByNameAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByNameAndHost(name, host string) ([]Record, error) {
	whereClause := "name = ? AND host = ?"
	return db.internalSelect(whereClause, name, host)
}

// SelectByUIDNameAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDNameAndHost(uid, name, host string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND host = ?"
	return db.internalSelect(whereClause, uid, name, host)
}

// SelectByUIDAndName prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDAndName(uid, name string) ([]Record, error) {
	whereClause := "uid = ? AND name = ?"
	return db.internalSelect(whereClause, uid, name)
}

// SelectByUID prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUID(uid string) ([]Record, error) {
	whereClause := "uid = ?"
	return db.internalSelect(whereClause, uid)
}

// internalSelect executes SQL SELECT query with given WHERE clause and fill the result(s) to array of the Record type.
func (db *Database) internalSelect(whereClause string, params ...interface{}) ([]Record, error) {
	query := "SELECT uid, name, host, data_type, data FROM migrations"
	if len(whereClause) > 0 {
		query += " WHERE " + whereClause
	}
	rows, err := db.conn.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SELECT query: %v", err)
	}
	defer rows.Close()
	selectResult := []Record{}
	for rows.Next() {
		rec := Record{}
		err = rows.Scan(&rec.UID, &rec.Name, &rec.Host, &rec.DataType, &rec.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		selectResult = append(selectResult, rec)
	}
	return selectResult, nil
}

// SelectAll executes SQL SELECT query and returns all records in the migration table.
func (db *Database) SelectAll() ([]Record, error) {
	return db.internalSelect("")
}

// Insert executes SQL INSERT query with a given record.
func (db *Database) Insert(rec []interface{}) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to start SQL transaction: %v", err)
	}
	sqlstmt := `INSERT INTO migrations(uid, name, host, data_type, data) values (?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(sqlstmt)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	_, err = stmt.Exec(rec...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}
	return tx.Commit()
}
