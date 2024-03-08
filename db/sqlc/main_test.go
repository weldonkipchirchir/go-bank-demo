package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
)

/*testQueries is a variable declared to hold an instance of the Queries struct, which encapsulates database queries and transactions. It's initialized later in the code.
 */
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	//function is called to establish a connection to the PostgreSQL database
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to db")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
