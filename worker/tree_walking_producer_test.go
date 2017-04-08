package worker

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestProducerWithFileFilter(t *testing.T) {
	tempDir := makeTempDir()
	defer cleanupTempDir(tempDir)
	jobsChan := make(chan Job, 10)

	// run Produce
	producer := TreeWalkingProducer{Root: tempDir, JobBuilder: jobBuilder, PathFilter: FileFilter}
	queued := producer.Produce(jobsChan)

	// Should get the right number of jobs queued
	if queued != 5 {
		t.Errorf("Wrong number of jobs queued (%d)\n", queued)
	}
	if len(jobsChan) != 5 {
		t.Errorf("Wrong number of jobs in channel (%d)\n", queued)
	}
}

func TestProduceWithNoFilter(t *testing.T) {
	tempDir := makeTempDir()
	defer cleanupTempDir(tempDir)
	jobsChan := make(chan Job, 10)

	// run Produce
	producer := TreeWalkingProducer{Root: tempDir, JobBuilder: jobBuilder}
	queued := producer.Produce(jobsChan)

	// Should get the right number of jobs queued (all files + two dirs
	if queued != 7 {
		t.Errorf("Wrong number of jobs queued (%d)\n", queued)
	}
}

// Fakes
type dummyJob struct {
}

func (job dummyJob) Work(chan Result) {
}

func jobBuilder(path string, info os.FileInfo, i int) Job {
	return dummyJob{}
}

// Make some temporary
func makeTempDir() (tempDir string) {
	// Make a temp directory with files and directories
	tempDir, _ = ioutil.TempDir("", "test-")
	for i := 0; i < 5; i++ {
		tempFile, _ := ioutil.TempFile(tempDir, "test-")
		tempFile.Close()
	}
	// make a sub directory
	ioutil.TempDir(tempDir, "test-")
	return tempDir
}

func cleanupTempDir(path string) {
	os.RemoveAll(path)
}
