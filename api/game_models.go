package api

type GameWithFriendsRequest struct {
	OponentId int32  `json:"oponent_id"`
	Subject   string `json:"subject"`
	Status    string `json:"status"`
	Grade     int32  `json:"grade"`
}

type GameWithRandomRequest struct {
	Subject string `json:"subject"`
}

type ConnectionResponse struct {
	Type             string `json:"type"`
	ConnectionStatus string `json:"connection_status"`
	OPName           string `json:"op_name"`
	OwnName          string `json:"own_name"`
}

type OponentIdResponse struct {
	OponentId int32 `json:"oponent_id"`
}
type GameQuestionsResponse struct {
	Type      string `json:"type"`
	ID        int32  `json:"id"`
	Questions string `json:"questions"`
	OptionsA  string `json:"options_a"`
	OptionsB  string `json:"options_b"`
	OptionsC  string `json:"options_c"`
	OptionsD  string `json:"options_d"`
	Correct   string `json:"correct"`
}
type GameAnswerRequest struct {
	CorrectAnswer string `json:"correct_answer"`
}
type CorrectResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Score   int32  `json:"score"`
	Oscore  int32  `json:"oscore"`
}
type WrongResponse struct {
	Type          string `json:"type"`
	Message       string `json:"message"`
	CorrectAnswer string `json:"correct_answer"`
	Score         int32  `json:"score"`
	Oscore        int32  `json:"oscore"`
}
type FinalResponse struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Score    int32  `json:"score"`
	Oscore   int32  `json:"oscore"`
	YourName string `json:"your_name"`
	OPName   string `json:"op_name"`
}
type Errors struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Error      error  `json:"error"`
	StatusCode int32  `json:"status_code"`
}
