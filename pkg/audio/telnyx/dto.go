package telnyx

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TranscriptionData struct {
	Confidence float64 `json:"confidence"`
	IsFinal    bool    `json:"is_final"`
	Transcript string  `json:"transcript"`
}

type RecordingUrls struct {
	Mp3 string `json:"mp3"`
	Wav string `json:"wav"`
}

type Event struct {
	Data struct {
		EventType  string `json:"event_type"`
		ID         string `json:"id"`
		OccurredAt string `json:"occurred_at"`
		Payload    struct {
			CallControlID     string            `json:"call_control_id"`
			CallLegID         string            `json:"call_leg_id"`
			CallSessionID     string            `json:"call_session_id"`
			ClientState       string            `json:"client_state"`
			ConnectionID      string            `json:"connection_id"`
			TranscriptionData TranscriptionData `json:"transcription_data"`
			CustomHeaders     []Header          `json:"custom_headers"`
			Direction         string            `json:"direction"`
			From              string            `json:"from"`
			State             string            `json:"state"`
			StartTime         string            `json:"start_time"`
			To                string            `json:"to"`
			// recording.error
			Reason string `json:"reason"`
			// recording.saved
			RecordingStartedAt  string        `json:"recording_started_at"`
			RecordingEndedAt    string        `json:"recording_ended_at"`
			Channels            string        `json:"channels"`              // single or dual
			RecordingUrls       RecordingUrls `json:"recording_urls"`        // Only valid for 10 minutes, unsure if we will use these
			PublicRecordingUrls RecordingUrls `json:"public_recording_urls"` // only if activated on app
		} `json:"payload"`
		// call.recording.saved
		RecordType  string `json:"record_type"`
	} `json:"data"`
	Meta struct {
		Attempt     int    `json:"attempt"`
		DeliveredTo string `json:"delivered_to"`
	} `json:"meta"`
}

type SoundModifications struct {
	Pitch    float64 `json:"pitch"`
	Semitone float64 `json:"semitone"`
	Octaves  float64 `json:"octaves"`
	Track    string  `json:"track"`
}

type CredentialsConfiguration struct {
	Bucket      string `json:"bucket"`
	Credentials string `json:"credentials"`
}

// Describes a recording as returned by the Telnyx Recording API
type Recording struct {
	CallControlID      string        `json:"call_control_id"`
	CallLegID          string        `json:"call_leg_id"`
	CallSessionID      string        `json:"call_session_id"`
	Channels           string        `json:"channels"` // single or dual
	ConferenceId       string        `json:"conference_id"`
	CreatedAt          string        `json:"created_at"`
	DownloadUrls       RecordingUrls `json:"download_urls"`
	DurationMillis     int           `json:"duration_millis"`
	ID                 string        `json:"id"`
	RecordType         string        `json:"record_type"` // "recording"
	RecordingStartedAt string        `json:"recording_started_at"`
	Source             string        `json:"source"` // "call" (or "conference" if conference call)
	Status             string        `json:"status"` // "completed"
	UpdatedAt          string        `json:"updated_at"`
}

type RecordingResponse struct {
	Data *Recording `json:"data"`
}
