package main

import (
	"fmt"

	"github.com/wtran29/streamer"
)

func main() {
	// define number of workers and jobs
	const numJobs = 4
	const numWorkers = 2

	// create chan for work and results
	notifyChan := make(chan streamer.ProcessingMessage, numJobs)
	defer close(notifyChan)

	videoQueue := make(chan streamer.VideoProcessingJob, numJobs)
	defer close(videoQueue)

	// get worker pool
	wp := streamer.New(videoQueue, numWorkers)
	fmt.Println("wp:", wp)

	// start worker pool
	wp.Run()

	// create 4 videos to send to worker pool

	// send videos to worker pool

	// print out results
}
