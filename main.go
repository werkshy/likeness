package main

import (
	"flag"
	"log"
	"os"

	"github.com/werkshy/likeness/index"
)

// Global config variables
var mainDir string

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.StringVar(&mainDir, "main-dir", "/data/delete_me/photos", "Main Photo Directory")

	//var dbHost = flag.String("db-host", "localhost", "Postgres DB Hostname")
	// TODO: add port, user, pass
	//var thumbDir = flag.String("thumb-dir", "/.thumbs", "Relative Path of Thumbnail storage")

	// TODO: add http server flags
	flag.Parse()

	var command = os.ExpandEnv(flag.Arg(0))

	checkConfig(mainDir)

	switch command {
	case "index":
		index.StartIndex(mainDir)
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
