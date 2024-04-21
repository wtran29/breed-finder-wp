package streamer

// ProcessingMessage is the info we send back to the client
type ProcessingMessage struct {
	ID         int
	Successful bool
	Message    string
	OutputFile string
}

// VideoProcessingJob is the unit of work performed.
// wraps around Video to have the info we need about input/output
type VideoProcessingJob struct {
	Video Video
}

type Processor struct {
	Engine Encoder
}

// Video is the type for a video that we wish to process
type Video struct {
	ID int
	InputFile string
	OutputDir string // pathname to where we want video to show up
	EncodingType string
	NotifyChan chan ProcessingMessage
	// Options *VideoOpts
	Encoder Processor
}
func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	workerPool := make(chan chan VideoProcessingJob, maxWorkers)

	// TODO implement processor logic
	p := Processor{}
	return &VideoDispatcher{
		jobQueue: jobQueue,
		maxWorkers: maxWorkers,
		WorkerPool: workerPool,
		Processor: p,
	}


}