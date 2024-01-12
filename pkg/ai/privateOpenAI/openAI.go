package privateOpenAI

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type OpenAIHandler struct {
	client         *azopenai.Client
	deploymentName string
}

func NewOpenAIHandler(apiKey, endpoint, deploymentName string) (*OpenAIHandler, error) {
	keyCredential := azcore.NewKeyCredential(apiKey)
	var options *azopenai.ClientOptions
	client, err := azopenai.NewClientWithKeyCredential(endpoint, keyCredential, options)
	if err != nil {
		log.Printf("Error creating OpenAI client: %v", err)
		return nil, err
	}

	return &OpenAIHandler{
		client:         client,
		deploymentName: deploymentName,
	}, nil
}

func (h *OpenAIHandler) GetSimpleChatCompletion(system string, text string) (string, error) {
	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: &system},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(text)},
	}

	resp, err := h.client.GetChatCompletions(
		context.Background(),
		azopenai.ChatCompletionsOptions{
			Messages:       messages,
			DeploymentName: &h.deploymentName,
		},
		nil,
	)
	if err != nil {
		log.Printf("Error getting chat completion from OpenAI: %v", err)
		return "", err
	}

	return *resp.Choices[0].Message.Content, nil
}
