package telnyx

type CredentialsPayload struct {
	Backend       string                   `json:"backend,omitempty"`
	Configuration CredentialsConfiguration `json:"configuration,omitempty"`
}

type AnswerPayload struct {
	BillingGroupID      string              `json:"billing_group_id,omitempty"`
	ClientState         string              `json:"client_state,omitempty"`
	CommandID           string              `json:"command_id,omitempty"`
	CustomHeaders       []Header            `json:"custom_headers,omitempty"`
	PreferredCodecs     string              `json:"preferred_codecs,omitempty"`
	SipHeaders          []Header            `json:"sip_headers,omitempty"`
	SoundModifications  *SoundModifications `json:"sound_modifications,omitempty"`
	StreamUrl           string              `json:"stream_url,omitempty"`
	StreamTrack         string              `json:"stream_track,omitempty"`           // "inbound_track", "outbound_track", "both_tracks"
	SendSilenceWhenIdle bool                `json:"send_silence_when_idle,omitempty"` // default false
	WebhookUrl          string              `json:"webhook_url,omitempty"`
	WebhookUrlMethod    string              `json:"webhook_url_method,omitempty"` // "POST", "GET"
}

type UpdateClientStatePayload struct {
	ClientState string `json:"client_state,omitempty"`
}

type GatherPayload struct {
	MinimumDigits           int    `json:"minimum_digits,omitempty"`             // default 1
	MaximumDigits           int    `json:"maximum_digits,omitempty"`             // default 128
	TimeoutMillis           int    `json:"timeout_millis,omitempty"`             // default 60_000
	InterDigitTimeoutMillis int    `json:"inter_digit_timeout_millis,omitempty"` // default 5_000
	InitialTimeoutMillis    int    `json:"initial_timeout_millis,omitempty"`     // default 5_000
	TerminatingDigit        string `json:"terminating_digit,omitempty"`          // default #
	ValidDigits             string `json:"valid_digits,omitempty"`               // default 0123456789*#
	GatherID                string `json:"gather_id,omitempty"`                  // randomly generated if not provided
	ClientState             string `json:"client_state,omitempty"`
	CommandID               string `json:"command_id,omitempty"`
}

type SimplePayload struct {
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
}

type RecordStartPayload struct {
	Format      string `json:"format,omitempty"`   // required either mp3 or wav
	Channels    string `json:"channels,omitempty"` // required either single or dual
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
	PlayBeep    bool   `json:"play_beep,omitempty"`
	MaxLength   int    `json:"max_length,omitempty"`   // default 0, max 14_400
	TimeoutSecs int    `json:"timeout_secs,omitempty"` // default 0 (infinite)
	Trim        string `json:"trim,omitempty"`         // default nil, other option: trim-silence
}

type SpeakTextPayload struct {
	Payload      string `json:"payload,omitempty"`       // Required string of max 3_000 characters
	PayloadType  string `json:"payload_type,omitempty"`  // Default text, other option: ssml
	ServiceLevel string `json:"service_level,omitempty"` // Default premium, other option: standard
	Stop         string `json:"stop,omitempty"`          // Default undefined, other option: current, all
	Voice        string `json:"voice,omitempty"`         // Required (Male or female)
	Language     string `json:"language,omitempty"`      // Required (nl-NL)
	ClientState  string `json:"client_state,omitempty"`
	CommandID    string `json:"command_id,omitempty"`
}

type NoiseSuppressionPayload struct {
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
	Direction   string `json:"direction,omitempty"` // Default inbound,  options: inbound, outbound, both
}

type TranscriptionPayload struct {
	Language            string `json:"language,omitempty"`        // Default "en", we'll use 'nl' or 'auto_detect'
	InterimResults      bool   `json:"interim_results,omitempty"` // Default false (only relevant for engine A)
	ClientState         string `json:"client_state,omitempty"`
	TranscriptionEngine string `json:"transcription_engine,omitempty"` // Default "A" (google), other option: "B" (telnyx)
	CommandID           string `json:"command_id,omitempty"`
}

type PlayAudio struct {
	AudioUrl        string `json:"audio_url,omitempty"`  // Cannot use MediaName and AudioUrl at the same time
	MediaName       string `json:"media_name,omitempty"` // Cannot use MediaName and AudioUrl at the same time (MediaName is the name of the audio file uploaded to Telnyx)
	Loop            int    `json:"loop,omitempty"`
	Overlay         bool   `json:"overlay,omitempty"`
	Stop            string `json:"stop,omitempty"`             // can be "current" or "all"
	TargetLegs      string `json:"target_legs,omitempty"`      // can be "self", "opposite" or "both"
	CacheAudio      bool   `json:"cache_audio,omitempty"`      // default true
	PlaybackContent string `json:"playback_content,omitempty"` // optional base64 encoded audio mp3 file.
	ClientState     string `json:"client_state,omitempty"`
	CommandID       string `json:"command_id,omitempty"`
}
