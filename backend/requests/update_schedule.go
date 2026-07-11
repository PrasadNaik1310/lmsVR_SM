package requests

type UpdateScheduleRequest struct {
	PlannedDate string `json:"planned_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	TeacherID   string `json:"teacher_id"`
	Status      string `json:"status"`
}
