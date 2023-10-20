package conversation

type Prompt struct{}

type UserType string

const (
	User   UserType = "user"
	Agent  UserType = "agent"
	System UserType = "system"
)

var validUserTypes = []UserType{User, Agent, System}

type ConversationStep struct {
	UserType UserType `json:"userType"`
	Text     string   `json:"text"`
	Prompt   *Prompt  `json:"prompt"`
	Purpose  string   `json:"purpose"`
}

type ClientState struct {
	Index   int    `json:"index"`
	Purpose string `json:"purpose"`
}

type ConversationRuleSet struct {
	ID     string `json:"id"`
	Simple bool   `json:"simple"`

	Steps []ConversationStep `json:"steps"`
}
