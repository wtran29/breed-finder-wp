package streamer

import "fmt"

// VideoDispatcher holds info for a dispatcher
type VideoDispatcher struct {
	WorkerPool chan chan VideoProcessingJob
	maxWorkers int
	jobQueue   chan VideoProcessingJob
	Processor  Processor
}

// videoWorker holds info for the pool worker.
type videoWorker struct {
	id         int
	jobQueue   chan VideoProcessingJob
	workerPool chan chan VideoProcessingJob
}

// newVideoWorker takes in an id and channel of chan VideoProcessingJob, returns a videoWorker
func newVideoWorker(id int, workerPool chan chan VideoProcessingJob) videoWorker {
	fmt.Println("newVideoWorker: creating video worker id", id)
	return videoWorker{
		id:         id,
		jobQueue:   make(chan VideoProcessingJob),
		workerPool: workerPool,
	}
}

// start() starts a worker
func (w videoWorker) start() {
	fmt.Println("w.start(): starting worker id", w.id)
	go func() {
		for {
			// add jobQueue to the worker pool
			w.workerPool <- w.jobQueue
			// wait for a job to come back
			job := <-w.jobQueue

			// process the job
			w.processVideoJob(job.Video)
		}
	}()
}

// Run() that will start everything
func (vd *VideoDispatcher) Run() {
	fmt.Println("vd.Run: starting worker pool by running workers")
	for i := 0; i < vd.maxWorkers; i++ {
		fmt.Println("vd.Run: starting worker id", i+1)
		worker := newVideoWorker(i+1, vd.WorkerPool)
		worker.start()
	}

	go vd.dispatch()
}

// dispatch() assign worker a job
func (vd *VideoDispatcher) dispatch() {
	for {
		// wait for job to come in
		job := <-vd.jobQueue
		fmt.Println("vd.dispatch: sending job", job.Video.ID, "to worker job queue")
		go func() {
			workerJobQueue := <-vd.WorkerPool
			workerJobQueue <- job
		}()
	}
}

// processingVideoJob
func (w videoWorker) processVideoJob(video Video) {
	fmt.Println("w.processVideoJob: starting encode on video", video.ID)
	video.encode()
}
