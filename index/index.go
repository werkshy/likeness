package index

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
)

type Result interface {
	Success() bool
	Get() interface{}
}

type Job interface {
	Work(results chan Result)
}

type Producer interface {
	Produce(jobs chan Job, done chan int)
}

type Consumer interface {
	Consume(jobs chan Job, results chan Result)
}

func ProduceConsume(producer Producer, consumer Consumer) {
	var bufferSize = 1000
	var jobs = make(chan Job, bufferSize)
	var doneProducing = make(chan int)
	var results = make(chan Result, bufferSize)

	go producer.Produce(jobs, doneProducing)

	go consumer.Consume(jobs, results)

	numQueued := <-doneProducing
	log.Printf("Queuing complete\n")

	for i := 0; i < numQueued; i++ {
		result := <-results
		switch value := result.Get().(type) {
		default:
			log.Printf("Unexpected type %T\n", value)
		case int:
			log.Printf("Job %d completed with success = %t\n", value, result.Success())

		}
	}
	log.Printf("%d workers complete\n", numQueued)

	//return e
}

// StartIndex kicks off an index job. Will wait until it finishes.
func StartIndex(mainDir string) (e error) {
	//log.Printf("Starting index of %s\n", mainDir)

	producer := TreeWalkingProducer{root: "/music"}
	consumer := SimpleConsumer{}

	ProduceConsume(producer, consumer)

	return e
}

type TreeWalkingProducer struct {
	root string
}

func (p TreeWalkingProducer) Produce(jobs chan Job, done chan int) {
	numQueued := 10
	for i := 0; i < numQueued; i++ {
		job := SleepyJob{id: i}
		jobs <- job
	}
	log.Printf("Done queuing jobs\n")
	done <- numQueued
}

type SimpleConsumer struct {
}

func (c SimpleConsumer) Consume(jobs chan Job, results chan Result) {
	for {
		// Pull out job
		job := <-jobs
		go job.Work(results)
	}
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
	var random, _ = rand.Int(rand.Reader, big.NewInt(1000))
	var sleepTime = random.Int64()
	fmt.Printf("Starting job %d, delay = %d\n", job.id, sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	job.success = true
	fmt.Printf("Finished job %d\n", job.id)
	results <- job
}
