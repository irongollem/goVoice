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

	Steps []ConversationStep `json:"steps" firestore:"steps"`
}

type ConversationStep struct {
	UserType string  `json:"userType" firestore:"userType"`
	Text     string  `json:"text" firestore:"text"`
	Prompt   *Prompt `json:"prompt" firestore:"prompt"`
	Purpose  string  `json:"purpose" firestore:"purpose"`
	AudioURL string  `json:"audioUrl" firestore:"audioUrl"`
}

type ConversationStepResponse struct {
	Purpose  string `json:"purpose" firestore:"purpose"`
	Response string `json:"response" firestore:"response"`
}

type Conversation struct {
	// conversation.ID should always be the same as the CallControlId
	ID               string            `firestore:"id"`
	RulesetID        string            `firestore:"rulesetId"`
	Responses        map[string]string `firestore:"responses"`
	Recordings       []Recording       `firestore:"recordings"`
	ConversationDone bool              `firestore:"conversationDone"`
}

/* ClientState with telnyx is a freeform, base64 encoded string to pass back and forth
 * Between the telnyx API and our application. It is used to keep track of
 * any state we want. Feel free to adjust this as needed.
 */
type ClientState struct {
	RulesetID        string `json:"rulesetId"`
	Purpose          string `json:"purpose"`
	CurrentStep      int    `json:"currentStep"`
	RecordingCount   int    `json:"recordingCount"`
	RecordingPurpose string `json:"recordingState"`
}

type Recording struct {
	Url     string `json:"url" firestore:"url"`
	Purpose string `json:"purpose" firestore:"purpose"`
}
