package streamer

import (
	"fmt"

	"github.com/xfrr/goffmpeg/transcoder"
)

// Encoder is an interface for encoding video. Any type that wants to satisfy
// this interface must implement all its methods.
type Encoder interface {
	EncodeToMP4(v *Video, baseFileName string) error
}

// VideoEncoder is a type that satisfies the Encoder interface by implementing
// all methods specied in Encoder
type VideoEncoder struct{}

// EncodeToMP4 takes a Video object and base filename and encodes to MP4 format
func (ve *VideoEncoder) EncodeToMP4(v *Video, baseFileName string) error {
	// create a transcoder
	trans := new(transcoder.Transcoder)

	// build output path
	outputPath := fmt.Sprintf("%s/%s", v.OutputDir, baseFileName)

	// initialize the transcoder
	err := trans.Initialize(v.InputFile, outputPath)
	if err != nil {
		return err
	}

	// set codec 
	trans.MediaFile().SetVideoCodec("libx264")

	// start transcoding process
	done := trans.Run(false)

	err = <-done
	if err != nil {
		return err
	}
	return nil
}
