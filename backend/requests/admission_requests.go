package requests

type CreateEnquiryRequest struct {
	FullName           string `json:"full_name" binding:"required"`
	Email              string `json:"email"`
	Phone              string `json:"phone" binding:"required"`
	InterestedCourseID string `json:"interested_course_id" binding:"required,uuid"`
	Notes              string `json:"notes"`
}

type UpdateEnquiryStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateApplicationRequest struct {
	EnquiryID       string `json:"enquiry_id" binding:"required,uuid"`
	AppliedCourseID string `json:"applied_course_id" binding:"required,uuid"`
}

type ApproveApplicationRequest struct {
	//TemporaryPassword string `json:"temporary_password"`
	ApplicationID string `json:"application_id" binding:"required,uuid"`
}

type RejectApplicationRequest struct {
	//Remarks string `json:"remarks"`
}
