package index

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"github.com/werkshy/likeness/worker"
)

// StartIndex kicks off an index job. Will wait until it finishes.
func StartIndex(mainDir string) (e error) {
	//log.Printf("Starting index of %s\n", mainDir)

	producer := TreeWalkingProducer{root: "/music"}
	consumer := worker.SimpleConsumer{}

	results := worker.ProduceConsume(producer, consumer)
	for _, result := range results {
		switch value := result.Get().(type) {
		//default:
		//	log.Printf("Unexpected type %T\n", value)
		case int:
			log.Printf("Job %d success = %t\n", value, result.Success())
		}
	}

	return e
}

type TreeWalkingProducer struct {
	root string
}

func (p TreeWalkingProducer) Produce(jobs chan worker.Job, done chan int) {
	numQueued := 10
	for i := 0; i < numQueued; i++ {
		job := SleepyJob{id: i}
		jobs <- job
	}
	log.Printf("Done queuing jobs\n")
	done <- numQueued
}

// Implement Job and Result
type SleepyJob struct {
	id      int
	success bool
}

func (job SleepyJob) Success() bool {
	return job.success
}

func (job SleepyJob) Get() interface{} {
	return job.id
}

func (job SleepyJob) Work(results chan worker.Result) {
	var random, _ = rand.Int(rand.Reader, big.NewInt(1000))
	var sleepTime = random.Int64()
	log.Printf("Starting job %d, delay = %d\n", job.id, sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	job.success = true
	log.Printf("Finished job %d\n", job.id)
	results <- job
}
