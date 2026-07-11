package requests

import "github.com/google/uuid"

type CreateScheduleRequest struct {
	//Title string `json:"title" binding:"required"`
	//Description string `json:"description"`
	LessonID    uuid.UUID `json:"lesson_id" binding:"required"`
	TeacherID   uuid.UUID `json:"teacher_id" binding:"required"`
	PlannedDate string    `json:"planned_date" binding:"required"`
	StartTime   string    `json:"start_time" binding:"required"`
	EndTime     string    `json:"end_time" binding:"required"`
}
