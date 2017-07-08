package index

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/werkshy/likeness/worker"
)

// curry a database connection into the jobs we build
func mkNewFileIndexJob(store Store) func(string, os.FileInfo, int) worker.Job {
	return func(path string, info os.FileInfo, i int) worker.Job {
		return &FileIndexJob{Store: store, id: i, path: path}
	}
}

// Implement Job and Result
type FileIndexJob struct {
	Store
	id             int
	path           string
	md5            []byte
	duplicatePhoto *Photo
	success        bool
}

func (result FileIndexJob) Success() bool {
	return result.success
}

func (result FileIndexJob) Get() interface{} {
	return result
}

func (result FileIndexJob) GetMd5() string {
	return fmt.Sprintf("%x", result.md5)
}

func (result FileIndexJob) isDupe() bool {
	return result.duplicatePhoto != nil
}

func (job *FileIndexJob) Work(results chan worker.Result) {
	found, err := job.FindPhotoByPath(job.path)
	switch {
	case err == sql.ErrNoRows:
		job.whenPhotoDoesNotExist(results)
	case err != nil:
		log.Printf("Error: Failed to lookup path: %s\n", err)
		job.success = false
	default:
		log.Printf("File is already in DB at %s\n", found.Path)
		job.success = true
	}
	// When the photo already exists in the DB, let's do nothing for now.

	results <- job
}

func (job *FileIndexJob) whenPhotoDoesNotExist(results chan worker.Result) {
	log.Printf("Hashing %s\n", job.path)
	job.md5 = fileHash(job.path)
	log.Printf("%s: %s\n", job.GetMd5(), job.path)

	// try to insert this photo. If there is a hash duplicate in the DB already,
	// mark it as a duplicate
	photo := Photo{
		Path:     job.path,
		Md5:      job.md5,
		FileDate: time.Now(), // FIXME
	}
	err := job.InsertPhoto(photo)
	switch {
	case err == UniqueConstraintViolation:
		dupe, err := job.FindPhotoByMd5(job.md5)
		if err == nil {
			job.duplicatePhoto = &dupe
			job.success = true
		} else {
			log.Printf("Failed to find the duplicate for %s\n", photo)
			job.success = false
		}
	case err != nil:
		log.Fatalf("Failed to insert: %s\n", err)
		job.success = false
	}

	results <- job
}

func fileHash(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}

	return hash.Sum(nil)
}
