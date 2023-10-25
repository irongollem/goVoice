package models

type Prompt struct {
	Text string `json:"text"`
}

type ConversationRuleSet struct {
	ID     string `json:"id"`
	Simple bool   `json:"simple"`

	Steps         []ConversationStep `json:"steps"`
	Conversations []Conversation     `json:"conversations"`
}

type ConversationStep struct {
	UserType string  `json:"userType"`
	Text     string  `json:"text"`
	Prompt   *Prompt `json:"prompt"`
	Purpose  string  `json:"purpose"`
}

type ConversationStepResponse struct {
	Purpose  string `json:"purpose"`
	Response string `json:"response"`
}

type Conversation struct {
	ID        string                      `json:"id"`
	Responses []ConversationStepResponse `json:"responses"`
}
