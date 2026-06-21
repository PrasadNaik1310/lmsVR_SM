package requests

type UpdateModuleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}
