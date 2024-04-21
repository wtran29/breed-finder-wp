package streamer

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

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
	ID           int
	InputFile    string
	OutputDir    string // pathname to where we want video to show up
	EncodingType string
	NotifyChan   chan ProcessingMessage
	Options      *VideoOpts
	Encoder      Processor
}

type VideoOpts struct {
	RenameOutput    bool
	SegmentDuration int // mp4 and hls format; produces a lot of ts files - knowing how long segments will be
	MaxRate1080p    string
	MaxRate720p     string
	MaxRate420p     string
}

func (vd *VideoDispatcher) NewVideo(id int, input, output, encType string, notifyChan chan ProcessingMessage, opts *VideoOpts) Video {
	if opts == nil {
		opts = &VideoOpts{}
	}
	return Video{
		ID:           id,
		InputFile:    input,
		OutputDir:    output,
		EncodingType: encType,
		NotifyChan:   notifyChan,
		Encoder:      vd.Processor,
		Options:      opts,
	}
}
func (v *Video) encode() {
	var fileName string
	switch v.EncodingType {
	case "mp4":
		// encode video
		name, err := v.encodeToMP4()
		if err != nil {
			// send info to notifyChan
			v.sendToNotifyChan(false, "", fmt.Sprintf("encode failed for %d %s", v.ID, err.Error()))
			return
		}
		fileName = fmt.Sprintf("%s.mp4", name)

	default:
		v.sendToNotifyChan(false, "", fmt.Sprintf("error processing for %d: invalid encoding type", v.ID))
		return
	}
	v.sendToNotifyChan(true, fileName, fmt.Sprintf("video id %d processed and saved as %s", v.ID, fmt.Sprintf("%s/%s", v.OutputDir, fileName)))
}

func (v *Video) encodeToMP4() (string, error) {
	baseFileName := ""

	if !v.Options.RenameOutput {
		// Get the base file name
		b := path.Base(v.InputFile)
		baseFileName = strings.TrimSuffix(b, filepath.Ext(b))
	} else {
		// TODO: generate random file name
	}
	err := v.Encoder.Engine.EncodeToMP4(v, baseFileName)
	if err != nil {
		return "", err
	}
	return baseFileName, nil
}

func (v *Video) sendToNotifyChan(successful bool, fileName, message string) {
	v.NotifyChan <- ProcessingMessage{
		ID:         v.ID,
		Successful: successful,
		Message:    message,
		OutputFile: fileName,
	}
}

// New creates and returns a video dispatcher
func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	workerPool := make(chan chan VideoProcessingJob, maxWorkers)

	// TODO implement processor logic
	var e VideoEncoder
	p := Processor{
		Engine: &e,
	}

	return &VideoDispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		WorkerPool: workerPool,
		Processor:  p,
	}

}
