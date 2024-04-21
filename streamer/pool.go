package streamer

// VideoDispatcher holds info for a dispatcher
type VideoDispatcher struct {
	WorkerPool chan chan VideoProcessingJob
	maxWorkers int
	jobQueue   chan VideoProcessingJob
	Processor  Processor
}

// videoWorker
type videoWorker struct {
	id         int
	jobQueue   chan VideoProcessingJob
	workerPool chan chan VideoProcessingJob
}

// newVideoWorker
func newVideoWorker(id int, workerPool chan chan VideoProcessingJob) videoWorker {
	return videoWorker{
		id:         id,
		jobQueue:   make(chan VideoProcessingJob),
		workerPool: workerPool,
	}
}

// start() starts a worker
func (w videoWorker) start() {
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
	for i := 0; i < vd.maxWorkers; i++ {
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

		go func() {
			workerJobQueue := <-vd.WorkerPool
			workerJobQueue <- job
		}()
	}
}

// processingVideoJob
func (w videoWorker) processVideoJob(video Video) {
	video.encode()
}
