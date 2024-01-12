package openAI

import (
	"context"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIHandler struct {
	client *openai.Client
}

func NewOpenAIHandler(apiKey string) *OpenAIHandler {
	client := openai.NewClient(apiKey)

	return &OpenAIHandler{
		client: client,
	}
}

func (h *OpenAIHandler) GetSimpleChatCompletion(system string, text string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: system,
		},
		{
			Role: openai.ChatMessageRoleUser,
			Content: text,
		},
	}

	resp, err := h.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		log.Printf("Error getting chat completion from OpenAI: %v", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
