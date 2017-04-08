package index

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/werkshy/likeness/worker"
)

// StartIndex kicks off an index job. Will wait until it finishes.
func StartIndex(mainDir string) (e error) {
	//log.Printf("Starting index of %s\n", mainDir)

	producer := worker.TreeWalkingProducer{
		Root:       mainDir,
		JobBuilder: newFileImportJob,
		PathFilter: worker.FileFilter,
	}

	results := worker.ProduceConsume(producer)
	for _, result := range results {
		switch value := result.Get().(type) {
		case FileImportJob:
			log.Printf("path: %s; hash: %s\n", value.path, value.md5)
		}
	}

	return e
}

func newFileImportJob(path string, info os.FileInfo, i int) worker.Job {
	return &FileImportJob{id: i, path: path}
}

// Implement Job and Result
type FileImportJob struct {
	id      int
	path    string
	md5     string
	success bool
}

func (result FileImportJob) Success() bool {
	return result.success
}

func (result FileImportJob) Get() interface{} {
	return result
}

func (job *FileImportJob) Work(results chan worker.Result) {
	job.md5 = fileHash(job.path)
	log.Printf("Finished %s\n", job.path)
	results <- job
}

func fileHash(path string) string {
	log.Printf("Hashing %s\n", path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
