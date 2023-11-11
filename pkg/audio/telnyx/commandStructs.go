package telnyx

type CredentialsPayload struct {
	Backend       string                   `json:"backend"`
	Configuration CredentialsConfiguration `json:"configuration"`
}

type AnswerPayload struct {
	BillingGroupID      string              `json:"billing_group_id"`
	ClientState         string              `json:"client_state"`
	CommandID           string              `json:"command_id"`
	CustomHeaders       []Header            `json:"custom_headers"`
	PreferredCodecs     string              `json:"preferred_codecs"`
	SipHeaders          []Header            `json:"sip_headers"`
	SoundModifications  *SoundModifications `json:"sound_modifications"`
	StreamUrl           string              `json:"stream_url"`
	StreamTrack         string              `json:"stream_track"`           // "inbound_track", "outbound_track", "both_tracks"
	SendSilenceWhenIdle bool                `json:"send_silence_when_idle"` // default false
	WebhookUrl          string              `json:"webhook_url"`
	WebhookUrlMethod    string              `json:"webhook_url_method"` // "POST", "GET"
}

type UpdateClientStatePayload struct {
	ClientState string `json:"client_state"`
}

type GatherPayload struct {
	MinimumDigits           int    `json:"minimum_digits"`             // default 1
	MaximumDigits           int    `json:"maximum_digits"`             // default 128
	TimeoutMillis           int    `json:"timeout_millis"`             // default 60_000
	InterDigitTimeoutMillis int    `json:"inter_digit_timeout_millis"` // default 5_000
	InitialTimeoutMillis    int    `json:"initial_timeout_millis"`     // default 5_000
	TerminatingDigit        string `json:"terminating_digit"`          // default #
	ValidDigits             string `json:"valid_digits"`               // default 0123456789*#
	GatherID                string `json:"gather_id"`                  // randomly generated if not provided
	ClientState             string `json:"client_state"`
	CommandID               string `json:"command_id"`
}

type SimplePayload struct {
	ClientState string `json:"client_state"`
	CommandID   string `json:"command_id"`
}

type RecordStartPayload struct {
	Format      string `json:"format"`   // required either mp3 or wav
	Channels    string `json:"channels"` // required either single or dual
	ClientState string `json:"client_state"`
	CommandID   string `json:"command_id"`
	PlayBeep    bool   `json:"play_beep"`
	MaxLength   int    `json:"max_length"`   // default 0, max 14_400
	TimeoutSecs int    `json:"timeout_secs"` // default 0 (infinite)
	Trim        string `json:"trim"`         // default nil, other option: trim-silence
}

type SpeakTextPayload struct {
	Payload      string `json:"payload"`       // Required string of max 3_000 characters
	PayloadType  string `json:"payload_type"`  // Default text, other option: ssml
	ServiceLevel string `json:"service_level"` // Default premium, other option: standard
	Stop         string `json:"stop"`          // Default undefined, other option: current, all
	Voice        string `json:"voice"`         // Required (Male or female)
	Language     string `json:"language"`      // Required (nl-NL)
	ClientState  string `json:"client_state"`
	CommandID    string `json:"command_id"`
}

type NoiseSuppressionPayload struct {
	ClientState string `json:"client_state"`
	CommandID   string `json:"command_id"`
	Direction   string `json:"direction"` // Default inbound,  options: inbound, outbound, both
}

type TranscriptionPayload struct {
	Language            string `json:"language"`        // Default "en", we'll use 'nl' or 'auto_detect'
	InterimResults      bool   `json:"interim_results"` // Default false (only relevant for engine A)
	ClientState         string `json:"client_state"`
	TranscriptionEngine string `json:"transcription_engine"` // Default "A" (google), other option: "B" (telnyx)
	CommandID           string `json:"command_id"`
}
