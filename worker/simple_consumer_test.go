package worker

import (
	"log"
	"testing"
)

type DummyResult struct {
}

func (r DummyResult) Success() bool {
	log.Printf("Success is good\n")
	return true
}

func (r DummyResult) Get() interface{} {
	return 1
}

type DummyJob struct {
	WasCalled bool
}

func (job *DummyJob) Work(results chan Result) {
	log.Printf("Job is called\n")
	job.WasCalled = true
	results <- DummyResult{}
}

func TestSimpleConsumer(t *testing.T) {
	// it pulls off a job from the channel and calls Work(results) on it
	simpleConsumer := SimpleConsumer{}
	jobs := make(chan Job, 1)
	results := make(chan Result, 1)
	job := new(DummyJob)
	jobs <- job

	go simpleConsumer.Consume(jobs, results)

	result := <-results
	if !result.Success() {
		t.Error("Result not successful")
	}
	if !job.WasCalled {
		t.Error("Job was not called")
	}

}
