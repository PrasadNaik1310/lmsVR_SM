package requests

type CreateLessonRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	ContentType     string `json:"content_type"`
	DurationMinutes int    `json:"duration_minutes"`
}
