package main

import (
	"flag"
	"log"
	"os"

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
	flag.StringVar(&mainDir, "main-dir", "/data/delete_me/photos", "Main Photo Directory")

	var dbUrl = flag.String("db-url", "postgres://localhost:5432/likeness", "Postgres DB Url")
	//var thumbDir = flag.String("thumb-dir", "/.thumbs", "Relative Path of Thumbnail storage")

	// TODO: add http server flags
	flag.Parse()

	var command = os.ExpandEnv(flag.Arg(0))

	checkConfig(mainDir)

	switch command {
	case "index":
		index.StartIndex(mainDir)
	case "migrate":
		schema.Migrate(dbUrl)
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
