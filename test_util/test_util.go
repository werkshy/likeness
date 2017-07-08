package test_util

// Some basic test setup helpers

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/werkshy/likeness/schema"
)

const TestDbUrl string = "postgres://localhost:5432/likeness-test"

func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func MigrateDb() {
	schema.Migrate(Db())
}

func Db() (db *sqlx.DB) {
	db = sqlx.MustConnect("postgres", TestDbUrl)
	return
}
