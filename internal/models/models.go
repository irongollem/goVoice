package models

type Prompt struct {
	Text string `json:"text"`
}

type Client struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ConversationRuleSet struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Simple bool    `json:"simple"`
	Client *Client `json:"client"`

	Steps []ConversationStep `json:"steps"`
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
	// conversation.ID should always be the same as the CallControlId
	ID               string            `firestore:"id"`
	RulesetID        string            `firestore:"rulesetId"`
	Responses        map[string]string `firestore:"responses"`
	Recordings     	 []Recording			 `firestore:"recordings"`
	ConversationDone bool              `firestore:"conversationDone"`
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
	ID             string `json:"id"`
	Url            string `json:"url"`
	ConversationID string `json:"conversationId"`
}
