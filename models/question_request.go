package models

type CreateQuestionRequest struct {
	Title             string       `json:"title"`
	MandatoryToAnswer bool         `json:"mandatory_to_answer"`
	QuestionType      string       `json:"question_type"`
	Options           []OptionReq  `json:"options"`
	Tags              []TagRequest `json:"tags"`
}

type OptionReq struct {
	Label     string `json:"label"`
	IsCorrect bool   `json:"is_correct"`
	Score     int64  `json:"score"`
}
