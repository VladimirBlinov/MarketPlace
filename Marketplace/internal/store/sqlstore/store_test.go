package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5433 user=admin password=qaz dbname=MarketPlaceTest sslmode=disable"
	}

	os.Exit(m.Run())
}
