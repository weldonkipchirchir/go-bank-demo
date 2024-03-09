package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/weldonkipchirchir/simple_bank/util"
)

/*testQueries is a variable declared to hold an instance of the Queries struct, which encapsulates database queries and transactions. It's initialized later in the code.
 */
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to load config")
	}
	//function is called to establish a connection to the PostgreSQL database
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
