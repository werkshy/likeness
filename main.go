package main

import (
	"flag"
	"fmt"
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
	fmt.Println("Usage: likeness COMMAND [ARGS]")
	fmt.Println("COMMAND is one of:")
	fmt.Println("  index: index existing photo dir")
	fmt.Println("ARGS and their default values:")
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

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		usage()
	}
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
