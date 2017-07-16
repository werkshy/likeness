package worker

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"
)

type DummyProducer struct {
	queueSize int
}

func (p DummyProducer) Produce(jobs chan Job) int {
	for i := 0; i < p.queueSize; i++ {
		job := SleepyJob{id: i}
		jobs <- job
	}
	log.Printf("Done queuing jobs\n")
	return p.queueSize
}

// Implement Job and Result
type SleepyJob struct {
	id      int
	success bool
}

func (job SleepyJob) Success() bool {
	return job.success
}

func (job SleepyJob) String() string {
	return fmt.Sprintf("%03d", job.id)
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

	results := ProduceConsume(producer)
	for _, result := range results {
		switch value := result.Get().(type) {
		case int:
			log.Printf("Job %d success = %t\n", value, result.Success())
		}
	}
}

func TestProduceConsumer(t *testing.T) {
	queueSize := 100

	producer := DummyProducer{queueSize: queueSize}

	results := ProduceConsume(producer)

	if len(results) != queueSize {
		t.Errorf("Was expecting %d results\n", queueSize)
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
