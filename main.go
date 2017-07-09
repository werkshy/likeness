package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	"github.com/werkshy/likeness/index"
	"github.com/werkshy/likeness/schema"
)

// Global config variables
var mainDir string

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&mainDir, "main-dir", "/home/oneill/test-photos", "Main Photo Directory")

	var dbUrl = flag.String("db-url", "postgres://localhost:5432/likeness", "Postgres DB Url")
	//var thumbDir = flag.String("thumb-dir", "/.thumbs", "Relative Path of Thumbnail storage")

	// TODO: add http server flags
	flag.Parse()

	var command = os.ExpandEnv(flag.Arg(0))

	checkConfig(mainDir)

	db := sqlx.MustConnect("postgres", *dbUrl)
	schema.Migrate(db)

	switch command {
	case "index":
		index.StartIndex(mainDir, db)
	default:
		log.Fatalf("Unknown command: '%s'\n", command)
	}
}

func checkConfig(mainDir string) {
	if _, err := os.Stat(mainDir); os.IsNotExist(err) {
		log.Fatalf("Main directory '%s' does not exist\n", mainDir)
	} else {
		log.Printf("Main dir: %s\n", mainDir)
	}

}
