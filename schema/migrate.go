package schema

import (
	"log"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func Migrate(dbUrl *string) {
	migrator := newMigrator(dbUrl)
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
		log.Printf("Current DB version: %d\n", version)
	}
}

// TODO make a Rollback method

func newMigrator(dbUrl *string) (migrator *migrate.Migrate) {
	migrator, err := migrate.New("file://migrations", *dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	return migrator
}
