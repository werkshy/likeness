package worker

import (
	"crypto/rand"
	"log"
	"math/big"
	"testing"
	"time"
)

type DummyProducer struct {
}

func (p DummyProducer) Produce(jobs chan Job, done chan int) {
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

func (job SleepyJob) Work(results chan Result) {
	var random, _ = rand.Int(rand.Reader, big.NewInt(10))
	var sleepTime = random.Int64()
	log.Printf("Starting job %d, delay = %d\n", job.id, sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	job.success = true
	log.Printf("Finished job %d\n", job.id)
	results <- job
}

func Example() {
	producer := DummyProducer{}
	consumer := SimpleConsumer{}

	results := ProduceConsume(producer, consumer)
	for _, result := range results {
		switch value := result.Get().(type) {
		case int:
			log.Printf("Job %d success = %t\n", value, result.Success())
		}
	}
}

func TestProduceConsumer(t *testing.T) {
	producer := DummyProducer{}
	consumer := SimpleConsumer{}

	results := ProduceConsume(producer, consumer)

	if len(results) != 10 {
		t.Error("Was expecting 10 results")
	}
	for _, result := range results {
		switch value := result.Get().(type) {
		default:
			t.Errorf("Unexpected type in result: %T\n", value)
		case int:
			if !result.Success() {
				t.Error("job failed")
			}
			log.Printf("Job %d success = %t\n", value, result.Success())
		}
	}
}
