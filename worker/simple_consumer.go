package worker

type SimpleConsumer struct {
}

func (c SimpleConsumer) Consume(jobs chan Job, results chan Result) {
	for {
		// Pull out job
		job := <-jobs
		go job.Work(results)
	}
}
