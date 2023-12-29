package models

type Prompt struct {
	Text string `json:"text" firestore:"text"`
}

type Client struct {
	Name  string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
}

type ConversationRuleSet struct {
	ID     string  `json:"id" firestore:"id"`
	Title  string  `json:"title" firestore:"title"`
	Simple bool    `json:"simple" firestore:"simple"`
	Client *Client `json:"client" firestore:"client"`

	Steps         []ConversationStep `json:"steps" firestore:"steps"`
}

type ConversationStep struct {
	UserType string  `json:"userType" firestore:"userType"`
	Text     string  `json:"text" firestore:"text"`
	Prompt   *Prompt `json:"prompt" firestore:"prompt"`
	Purpose  string  `json:"purpose" firestore:"purpose"`
	AudioFile string `json:"audioFile" firestore:"audioFile"`
}

type ConversationStepResponse struct {
	Purpose  string `json:"purpose" firestore:"purpose"`
	Response string `json:"response" firestore:"response"`
}

type Conversation struct {
	// conversation.ID should always be the same as the CallControlId
	ID               string                     `firestore:"id"`
	RulesetID        string                     `firestore:"rulesetId"`
	Responses        []ConversationStepResponse `firestore:"responses"`
	Recording        *Recording                 `firestore:"recordings"`
	ConversationDone bool                       `firestore:"conversationDone"`
}

/* ClientState with telnyx is a freeform, base64 encoded string to pass back and forth
 * Between the telnyx API and our application. It is used to keep track of
 * any state we want. Feel free to adjust this as needed.
 */
type ClientState struct {
	RulesetID   string `json:"rulesetId"`
	CurrentStep int    `json:"currentStep"`
}

type Recording struct {
	ID             string `json:"id" firestore:"id"`
	Url            string `json:"url" firestore:"url"`
	ConversationID string `json:"conversationId" firestore:"conversationId"`
}
