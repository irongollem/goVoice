package audioProcessor

// The AudioProcessor handles the denoising of the incoming audio before it is sent to the transcriber.

type AudioProcessor struct {}

func (lp *AudioProcessor) Transcribe(arr []byte) (string, error) {
	return "", nil
}

