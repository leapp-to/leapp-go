package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

var testDB *Database

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	teardown()
	os.Exit(c)
}

func setup() {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to the DB: %v\n", err)
		os.Exit(1)
	}
	testDB = &Database{conn: conn}
	sql := `
	  CREATE TABLE IF NOT EXISTS migrations(
		uid TEXT,
		name TEXT,
		host TEXT,
		data_type TEXT,
		data TEXT,
		PRIMARY KEY(uid, name, host, data_type)
	);
	`
	if err := testDB.Exec(sql); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create table: %v\n", err)
		os.Exit(1)
	}
}

func teardown() {
	if err := testDB.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close DB connection: %v\n", err)
		os.Exit(1)
	}
}

func TestExec(t *testing.T) {
	succStmt := "SELECT * from migrations;"
	err := testDB.Exec(succStmt)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}
}

func TestSelects(t *testing.T) {
	rec := map[string]string{
		"UID":      "testUID",
		"name":     "test",
		"host":     "testHost.io",
		"dataType": "testType",
		"data":     "testData"}
	//SelectAll()
	_, err := testDB.SelectAll()
	checkErr(t, "SelectAll", err)

	//SelectByUID
	_, err = testDB.SelectByUID(rec["UID"])
	checkErr(t, "SelectByUID", err)
	//SelectByUIDAndName
	_, err = testDB.SelectByUIDAndName(rec["UID"], rec["name"])
	checkErr(t, "SelectByUIDAndName", err)
	//SelectByUIDNameAndHost
	_, err = testDB.SelectByUIDNameAndHost(rec["UID"], rec["name"], rec["host"])
	checkErr(t, "SelectByUIDNameAndHost", err)
	//SelectByNameAndHost
	_, err = testDB.SelectByNameAndHost(rec["name"], rec["host"])
	checkErr(t, "SelectByNameAndHost", err)
	//SelectByUIDNameAndType
	_, err = testDB.SelectByUIDNameAndType(rec["UID"], rec["name"], rec["dataType"])
	checkErr(t, "SelectByUIDNameAndType", err)
	//SelectByUIDNameTypeAndHost
	_, err = testDB.SelectByUIDNameTypeAndHost(rec["UID"], rec["name"], rec["dataType"], rec["host"])
	checkErr(t, "SelectByUIDNameTypeAndHost", err)
}

func checkErr(t *testing.T, funcName string, err error) {
	if err != nil {
		t.Errorf("%v failed: %v", funcName, err)
	}
}

func TestInsert(t *testing.T) {
	rec := []interface{}{"testUID", "test", "testHost.io", "testType", "testData"}
	err := testDB.Insert(rec)
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}
}
