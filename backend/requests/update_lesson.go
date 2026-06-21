package requests

type UpdateLessonRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	ContentType     string `json:"content_type"`
	DurationMinutes int    `json:"duration_minutes"`
	IsPublished     bool   `json:"is_published"`
}
