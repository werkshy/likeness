package schema

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// Store the absolute path of the current dir so we can find the migrations dir
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func Migrate(db *sqlx.DB) {
	start := time.Now()
	checkPostgresVersion(db)
	migrator := newMigrator(db.DB)
	version, dirty, err := migrator.Version()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Current version: %d  dirty=%t\n", version, dirty)
	err = migrator.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	version, _, err = migrator.Version()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Current schema version: %d\n", version)
	}
	log.Printf("Migration check took %s\n", time.Since(start))
}

// TODO make a Rollback method

func newMigrator(db *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	migrationsDir := fmt.Sprintf("file:///%s/../migrations", basepath)
	migrator, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	return migrator
}

// ON CONFLICT clause requires Postgres >= 9.5
func checkPostgresVersion(db *sqlx.DB) {
	var serverVersion string
	err := db.Get(&serverVersion, "SHOW server_version")
	if err != nil {
		panic(err)
	}

	if serverVersion <= "9.5.0" {
		log.Fatalf("Detected postgres version '%s', need at least 9.5\n", serverVersion)
	}

}
