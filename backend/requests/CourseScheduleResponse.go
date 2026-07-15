package requests

import (
	"time"

	"github.com/google/uuid"
)

type CourseScheduleResponse struct {
	ID               uuid.UUID `json:"id"`
	LessonID         uuid.UUID `json:"lesson_id"`
	LessonTitle      string    `json:"lesson_title"`
	TeacherID        uuid.UUID `json:"teacher_id"`
	TeacherName      string    `json:"teacher_name"`
	PlannedDate      time.Time `json:"planned_date"`
	PlannedStartTime string    `json:"planned_start_time"`
	PlannedEndTime   string    `json:"planned_end_time"`
	Status           string    `json:"status"`
}
