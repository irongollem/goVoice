package ai

import (
	"goVoice/internal/config"
	// "goVoice/pkg/ai/openAI"
	"goVoice/pkg/ai/privateOpenAI"
)

type AIProvider interface {
	GetSimpleChatCompletion(system string, text string) (string, error)
}

func InitiateAIProvider (cfg *config.Config) (AIProvider, error) {
	// provider := openAI.NewOpenAIHandler(cfg.OpenAIKey)
	return privateOpenAI.NewOpenAIHandler(cfg.OpenAIKey, cfg.OpenAIEndpoint, cfg.OpenAIDeploymentName)
}