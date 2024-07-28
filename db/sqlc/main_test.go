package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weldonkipchirchir/simple_bank/util"
)

/*testQueries is a variable declared to hold an instance of the Queries struct, which encapsulates database queries and transactions. It's initialized later in the code.
 */
var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to load config")
	}
	//function is called to establish a connection to the PostgreSQL database
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db")
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
