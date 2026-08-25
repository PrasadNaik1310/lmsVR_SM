package requests

type CreateAcademicSessionRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	IsActive  *bool  `json:"is_active"`
}

type UpdateAcademicSessionRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  *bool  `json:"is_active"`
}
