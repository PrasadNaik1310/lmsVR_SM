package requests

type CreateCourseRequest struct {
	Title             string `json:"title" binding:"required"`
	Description       string `json:"description"`
	Level             string `json:"level"`
	AcademicSessionID string `json:"academic_session_id" `
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	TotalSeats        int    `json:"total_seats"`
	MeetLink          string `json:"meet_link"`
}
