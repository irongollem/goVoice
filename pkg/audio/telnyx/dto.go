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

type Event struct {
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
	} `json:"payload"`
	RecordType string `json:"record_type"`
	Meta       struct {
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

type CommandPayload struct {
	CommandId          string              `json:"command_id"`
	CustomHeaders      []Header            `json:"custom_headers"`
	SipHeaders         []Header            `json:"sip_headers"`
	SoundModifications *SoundModifications `json:"sound_modifications"`
	ClientState        string              `json:"client_state"`
	// for transcription use nl for speak use nl-NL
	Language string `json:"language"`

	// Answer specific
	BillingGroupID      string `json:"billing_group_id"`
	StreamUrl           string `json:"stream_url"`
	StreamTrack         string `json:"stream_track"`
	SendSilenceWhenIdle bool   `json:"send_silence_when_idle"`
	WebhookUrl          string `json:"webhook_url"`
	WehookUrlMethod     string `json:"webhook_url_method"`

	// Transcription specific
	TranscriptionEngine string `json:"transcription_engine"`

	// Recording specific
	Format      string `json:"format"`
	Channels    string `json:"channels"`
	Trim        string `json:"trim"`
	PlayBeap    bool   `json:"play_beep"`
	MaxLength   int    `json:"max_length"`
	TimeoutSecs int    `json:"timeout_secs"`

	// Speak specific
	Payload      string `json:"payload"`
	PayloadType  string `json:"payload_type"`  // text or ssml
	ServiceLevel string `json:"service_level"` // premium or basic (default premium needed for non en-US)
	Stop         string `json:"stop"`          // undefined, current or all to stop any playing audio
	Voice        string `json:"voice"`         // male or female
}
