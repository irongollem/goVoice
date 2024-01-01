package ai

type ValidatedAnswer struct {
	Answer string `json:"answer"`
	Purpose string `json:"purpose"`
	Question string `json:"question"`
	Complete *bool `json:"complete"`
}