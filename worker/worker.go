package worker

import "log"

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

func ProduceConsume(producer Producer, consumer Consumer) (results []Result) {
	var bufferSize = 1000
	var jobs = make(chan Job, bufferSize)
	var doneProducing = make(chan int)
	var resultsChan = make(chan Result, bufferSize)

	go producer.Produce(jobs, doneProducing)

	go consumer.Consume(jobs, resultsChan)

	numQueued := <-doneProducing
	log.Printf("Queuing complete\n")

	for i := 0; i < numQueued; i++ {
		result := <-resultsChan
		results = append(results, result)
	}
	log.Printf("%d workers complete\n", numQueued)

	return results
}
