package mysql

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {

	// Create a new database connection
	testDB, err := sql.Open("mysql", "test_user:test_pass@/sniptext_test?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Setup the database tables
	query, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Run the queries for setup
	_, err = testDB.Exec(string(query))
	if err != nil {
		t.Fatal(err)
	}

	// Return the database connection and a function for teardown
	return testDB, func() {
		// Teardown the database tables
		query, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		// Run the queries for teardown
		_, err = testDB.Exec(string(query))
		if err != nil {
			t.Fatal(err)
		}

		testDB.Close()
	}
}
