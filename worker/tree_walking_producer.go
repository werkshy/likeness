package worker

import (
	"log"
	"os"
	"path/filepath"
)

// JobBuildFunc is the type passed into TreeWalkingProducer that instantiates a
// path-visiting Job.
type JobBuildFunc func(path string, info os.FileInfo, i int) Job

// PathFilterFunc is the type optionally passed into TreeWalkingProducer that
// decides whether or not this path needs a visiting job created.
type PathFilterFunc func(path string, info os.FileInfo) bool

type visitResult struct {
	path string
	info os.FileInfo
}

// TreeWalkingProducer visits every file in the filesystem and queues a Job for
// each file.
type TreeWalkingProducer struct {
	Root       string
	JobBuilder JobBuildFunc
	PathFilter PathFilterFunc
	results    []visitResult
}

// Produce (implement the Producer interface) - visit every path under root and enqueue a
// Job
func (p TreeWalkingProducer) Produce(jobs chan Job) int {
	queued := 0
	filepath.Walk(p.Root, p.visit)
	log.Printf("There are %d paths to consider\n", len(p.results))

	for i, result := range p.results {
		if p.shouldQueue(result) {
			job := p.JobBuilder(result.path, result.info, i)
			jobs <- job
			queued += 1
		}
	}
	log.Printf("Done queuing jobs\n")
	return queued
}

func (p TreeWalkingProducer) shouldQueue(result visitResult) bool {
	if p.PathFilter == nil {
		return true
	}
	return p.PathFilter(result.path, result.info)
}

// The WalkFunc for TreeWalkingProducer, simply builds a list of paths
func (p *TreeWalkingProducer) visit(path string, info os.FileInfo, err error) error {
	p.results = append(p.results, visitResult{path: path, info: info})

	return nil
}
