package cmd

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/werkshy/likeness/schema"
)

var mainDir string
var dbURL string
var dbConnection *sqlx.DB

// RootCmd implements main 'root' cobra arg parsing, for args that
// apply to all subcommands.
var RootCmd = &cobra.Command{
	Use:   "likeness",
	Short: "Likeness is a photo importer and organizer",
	Long:  `Import, organize, detect duplicated photo files.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkConfig()
		dbConnection = createDbConnection()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Yo\n")
	},
}

func init() {
	flags := RootCmd.PersistentFlags()
	flags.StringVar(&mainDir, "main-dir", "/home/oneill/test-photos", "Main Photo Directory")

	flags.StringVar(&dbURL, "db-url", "postgres://localhost:5432/likeness", "Postgres DB Url")
	//var thumbDir = flag.String("thumb-dir", "/.thumbs", "Relative Path of Thumbnail storage")
}

func checkConfig() {
	if _, err := os.Stat(mainDir); os.IsNotExist(err) {
		log.Fatalf("Main directory '%s' does not exist\n", mainDir)
	} else {
		log.Printf("Main dir: %s\n", mainDir)
	}
}

func createDbConnection() (db *sqlx.DB) {
	db = sqlx.MustConnect("postgres", dbURL)
	schema.Migrate(db)
	return
}
