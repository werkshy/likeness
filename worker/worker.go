package worker

import (
	"log"
	"runtime"
)

type Result interface {
	Success() bool
	Get() interface{}
}

type Job interface {
	Work(results chan Result)
	String() string
}

type Producer interface {
	Produce(jobs chan Job) int
}

func ProduceConsume(producer Producer) (results []Result) {
	numWorkers := runtime.NumCPU()

	var bufferSize = 100000
	var jobs = make(chan Job, bufferSize)
	var resultsChan = make(chan Result, bufferSize)

	for i := 0; i < numWorkers; i++ {
		go consume(jobs, resultsChan)
	}

	numQueued := producer.Produce(jobs)
	close(jobs)
	log.Printf("Queuing complete (%d)\n", numQueued)

	for i := 0; i < numQueued; i++ {
		result := <-resultsChan
		results = append(results, result)
	}
	log.Printf("%d workers complete, we have %d results\n", numQueued, len(results))

	return results
}

func consume(jobs chan Job, results chan Result) {
	for job := range jobs {
		job.Work(results)
	}
}
