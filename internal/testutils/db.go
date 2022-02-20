package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

// Schema used by tables.
const tablesSchema = "public"

// Tables to ignore when flushing the database.
var ignoreTables = []string{"schema_migrations"}

// PrepareDB creates a connection with the test database and flush the existing data.
func PrepareDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	err = flushDatabase(db)
	if err != nil {
		t.Fatal(fmt.Errorf("error during flushing database: %w", err))
	}

	return db
}

// Deletes all tables values in a single transaction
func flushDatabase(db *sql.DB) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	tables, err := getTablesName(db)
	if err != nil {
		log.Fatal(err)
	}

OuterLoop:
	for _, table := range tables {
		// Remove the schema.
		table = strings.Replace(table, tablesSchema+".", "", 1)

		// Ignore tables we don't need to flush.
		for _, t := range ignoreTables {
			if table == t {
				continue OuterLoop
			}
		}

		_, err = tx.ExecContext(ctx, fmt.Sprintf(`TRUNCATE "%s" CASCADE`, table))
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// getTablesName returns a list of the tables in the database.
func getTablesName(db *sql.DB) ([]string, error) {
	var tables []string

	const query = `
		SELECT pg_namespace.nspname || '.' || pg_class.relname
		FROM pg_class
		INNER JOIN pg_namespace ON pg_namespace.oid = pg_class.relnamespace
		WHERE pg_class.relkind = 'r'
		  AND pg_namespace.nspname NOT IN ('pg_catalog', 'information_schema', 'crdb_internal')
		  AND pg_namespace.nspname NOT LIKE 'pg_toast%'
		  AND pg_namespace.nspname NOT LIKE '\_timescaledb%';
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
