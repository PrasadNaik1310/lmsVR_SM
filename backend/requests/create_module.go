package requests

type CreateModuleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}
