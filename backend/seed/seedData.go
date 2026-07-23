package seed

import (
	"time"

	"github.com/google/uuid"

	"github.com/PrasadNaik1310/LMSVR_SM/models"
)

// Pre-defined UUIDs for referential integrity across Enquiries and Applications
var (
	// Course IDs (assumed to already exist in the DB)
	courseGoID     = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000001")
	coursePythonID = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000002")
	courseWebDevID = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000003")
	courseDataID   = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000004")

	// Enquiry IDs (used to link Applications)
	enquiry1ID = uuid.MustParse("bbbbbbbb-0000-0000-0000-000000000001")
	enquiry2ID = uuid.MustParse("bbbbbbbb-0000-0000-0000-000000000002")
	enquiry3ID = uuid.MustParse("bbbbbbbb-0000-0000-0000-000000000003")
	enquiry4ID = uuid.MustParse("bbbbbbbb-0000-0000-0000-000000000004")
	enquiry5ID = uuid.MustParse("bbbbbbbb-0000-0000-0000-000000000005")
)

// EnquirySeeds returns demo Enquiry records.
var EnquirySeeds = []models.Enquiry{
	{
		ID:                 enquiry1ID,
		FullName:           "Arjun Mehta",
		Email:              "arjun.mehta@example.com",
		Phone:              "+91-9393948568",
		InterestedCourseID: courseGoID,
		Status:             "new",
		Notes:              "Interested in weekend batches for Go programming.",
		CreatedAt:          time.Now().AddDate(0, 0, -10),
	},
	{
		ID:                 enquiry2ID,
		FullName:           "Priya Sharma",
		Email:              "priya.sharma@example.com",
		Phone:              "+91-0184937594",
		InterestedCourseID: coursePythonID,
		Status:             "contacted",
		Notes:              "Needs fee structure details. Prefers evening slots.",
		CreatedAt:          time.Now().AddDate(0, 0, -8),
	},
	{
		ID:                 enquiry3ID,
		FullName:           "Rohan Desai",
		Email:              "rohan.desai@example.com",
		Phone:              "+91-9322223345",
		InterestedCourseID: courseWebDevID,
		Status:             "follow_up",
		Notes:              "Asked about placement support after completion.",
		CreatedAt:          time.Now().AddDate(0, 0, -5),
	},
	{
		ID:                 enquiry4ID,
		FullName:           "Sneha Patel",
		Email:              "sneha.patel@example.com",
		Phone:              "+91-9001122334",
		InterestedCourseID: courseDataID,
		Status:             "converted",
		Notes:              "Ready to enroll. Waiting for admission form.",
		CreatedAt:          time.Now().AddDate(0, 0, -3),
	},
	{
		ID:                 enquiry5ID,
		FullName:           "Vikram Nair",
		Email:              "vikram.nair@example.com",
		Phone:              "+91-9765432109",
		InterestedCourseID: coursePythonID,
		Status:             "closed",
		Notes:              "Not interested at this time. May revisit next quarter.",
		CreatedAt:          time.Now().AddDate(0, 0, -1),
	},
}

// ApplicationSeeds returns demo Application records linked to the enquiries above.
var ApplicationSeeds = []models.Application{
	{
		ID:                uuid.New(),
		EnquiryID:         enquiry1ID,
		AppliedCourseID:   courseGoID,
		ApplicationStatus: "submitted",
		Remarks:           "All documents verified. Awaiting counsellor review.",
		SubmittedAt:       time.Now().AddDate(0, 0, -9),
	},
	{
		ID:                uuid.New(),
		EnquiryID:         enquiry2ID,
		AppliedCourseID:   coursePythonID,
		ApplicationStatus: "under_review",
		Remarks:           "Academic background check in progress.",
		SubmittedAt:       time.Now().AddDate(0, 0, -7),
	},
	{
		ID:                uuid.New(),
		EnquiryID:         enquiry3ID,
		AppliedCourseID:   courseWebDevID,
		ApplicationStatus: "approved",
		Remarks:           "Approved. Fee payment link sent via email.",
		SubmittedAt:       time.Now().AddDate(0, 0, -4),
	},
	{
		ID:                uuid.New(),
		EnquiryID:         enquiry4ID,
		AppliedCourseID:   courseDataID,
		ApplicationStatus: "enrolled",
		Remarks:           "Fee received. Student onboarded to LMS.",
		SubmittedAt:       time.Now().AddDate(0, 0, -2),
	},
	{
		ID:                uuid.New(),
		EnquiryID:         enquiry5ID,
		AppliedCourseID:   coursePythonID,
		ApplicationStatus: "rejected",
		Remarks:           "Applicant withdrew request before processing.",
		SubmittedAt:       time.Now().AddDate(0, 0, -1),
	},
}
