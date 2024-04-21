package main

import (
	"fmt"

	"github.com/wtran29/streamer"
)

func main() {
	// define number of workers and jobs
	const numJobs = 4
	const numWorkers = 1

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
	video := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "mp4", notifyChan, nil)
	// send videos to worker pool
	videoQueue <- streamer.VideoProcessingJob{Video: video}
	// print out results

}
