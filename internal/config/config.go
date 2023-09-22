package config

type Config struct {
	WebClientAPIAddr string
	VoiceAPIAddr     string
}

func LoadConfig() (*Config, error) {
	return &Config{
		WebClientAPIAddr: ":8080",
		VoiceAPIAddr:     ":8081",
	}, nil
}
