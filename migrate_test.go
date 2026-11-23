package sp_model

import (
	"os"
	"testing"
)

// TestCleanMigration runs the migrations of both databases against empty PostgreSQL
// databases. It is skipped unless SP_TEST_DATA_DSN and SP_TEST_INFO_DSN are set.
func TestCleanMigration(t *testing.T) {
	data, info := os.Getenv("SP_TEST_DATA_DSN"), os.Getenv("SP_TEST_INFO_DSN")
	if data == "" || info == "" {
		t.Skip("SP_TEST_DATA_DSN and SP_TEST_INFO_DSN not set")
	}
	ClientsInit(data, info) // panics when a migration fails
}
