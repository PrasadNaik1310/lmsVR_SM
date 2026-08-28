package requests

type CreateCompanyUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Phone     string `json:"phone"`

	Role string `json:"role" binding:"required"`

	// Student fields
	EnrollmentNumber string `json:"enrollment_number"`
	DateOfBirth      string `json:"date_of_birth"`
	Address          string `json:"address"`
	AdmissionDate    string `json:"admission_date"`

	// Teacher fields
	Specialization string `json:"specialization"`
	Bio            string `json:"bio"`
	JoiningDate    string `json:"joining_date"`
}
