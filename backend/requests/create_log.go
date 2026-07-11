package requests

type CreateCourseLogRequest struct {
	ConductedDate    string `json:"conducted_date" binding:"required"`
	CompletionStatus string `json:"completion_status" binding:"required"`
	Remarks          string `json:"remarks"`
	Homework         string `json:"homework"`
	NextTopic        string `json:"next_topic"`
}
