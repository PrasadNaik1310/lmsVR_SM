package requests

type UpdateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	TotalSeats  int    `json:"total_seats"`
	MeetLink    string `json:"meet_link"`
}
