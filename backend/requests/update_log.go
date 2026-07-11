package requests

type UpdateCourseLogRequest struct {
	CompletionStatus string `json:"completion_status"`
	Remarks          string `json:"remarks"`
	Homework         string `json:"homework"`
	NextTopic        string `json:"next_topic"`
}
