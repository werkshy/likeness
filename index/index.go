package index

// Index: go through the main photo dir and make sure we know all of the files
// within.
// Print out a list of duplicates

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/werkshy/likeness/worker"
)

// StartIndex kicks off an index job. Will wait until it finishes.
func StartIndex(mainDir string, db *sqlx.DB) (e error) {
	//log.Printf("Starting index of %s\n", mainDir)
	store := NewDbStore(db)

	producer := worker.TreeWalkingProducer{
		Root:       mainDir,
		JobBuilder: mkNewFileIndexJob(store),
		PathFilter: worker.FileFilter,
	}

	results := worker.ProduceConsume(producer)
	var duplicates []FileIndexJob

	for _, result := range results {
		switch job := result.Get().(type) {
		case FileIndexJob:
			log.Printf("[%s] %s\n", job.StatusString(), job.path)
			if job.isDupe() {
				duplicates = append(duplicates, job)
			}
		}
	}

	for _, job := range duplicates {
		fmt.Printf("%s is a duplicate of %s\n", job.path, job.duplicatePhoto.Path)
	}
	fmt.Printf("Found %d dupes\n", len(duplicates))

	return e
}
