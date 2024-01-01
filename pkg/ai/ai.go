package ai

import (
	"goVoice/internal/config"
	"goVoice/pkg/ai/openAI"
)

type AIProvider interface {
	GetSimpleChatCompletion(system string, text string) (string, error)
}

func InitiateAIProvider (cfg *config.Config) AIProvider {
	provider := openAI.NewOpenAIHandler(cfg.OpenAIKey)

	return provider
}