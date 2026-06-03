package handlers_test

import (
	"bytes"
	//"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"

	//"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestCreateEnquiry(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	courseID := uuid.New()

	// Course lookup
	mock.ExpectQuery(`SELECT .* FROM "courses" WHERE id = \$1 ORDER BY "courses"\."id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "level", "thumbnail_url", "status", "total_seats", "booked_seats", "start_date", "end_date", "meet_link", "created_by", "academic_session_id", "created_at", "updated_at"}).
			AddRow(courseID, "Backend Development", "", "", "", "published", 0, 0, nil, nil, "", uuid.Nil, uuid.Nil, nil, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO .*"enquiries"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	payload := []byte(`{
		"full_name":"Prasad Patil",
		"email":"prasad@example.com",
		"phone":"9876543210",
		"interested_course_id":"` + courseID.String() + `",
		"notes":"Interested in backend development"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/lms/admissions/enquiries", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	handlers.CreateEnquiry(c)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if response["enquiry"] == nil {
		t.Fatalf("expected enquiry in response")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListEnquiries(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	enquiryID := uuid.New()
	courseID := uuid.New()

	// Count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "enquiries"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// List query
	mock.ExpectQuery(`SELECT .* FROM "enquiries" LIMIT \$1`).
		WithArgs(20).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "interested_course_id", "status", "notes", "created_at"}).
			AddRow(enquiryID, "Prasad Patil", "prasad@example.com", "9876543210", courseID, "new", "Interested", time.Now()))

	req := httptest.NewRequest(http.MethodGet, "/lms/admissions/enquiries", nil)
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	handlers.ListEnquiries(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if response["enquiries"] == nil {
		t.Fatalf("expected enquiries in response")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetEnquiry(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	enquiryID := uuid.New()
	courseID := uuid.New()

	// Enquiry lookup
	mock.ExpectQuery(`SELECT .* FROM "enquiries" WHERE id = \$1 ORDER BY "enquiries"\."id" LIMIT \$2`).
		WithArgs(enquiryID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "interested_course_id", "status", "notes", "created_at"}).
			AddRow(enquiryID, "Prasad Patil", "prasad@example.com", "9876543210", courseID, "new", "Interested", time.Now()))

	req := httptest.NewRequest(http.MethodGet, "/lms/admissions/enquiries/"+enquiryID.String(), nil)
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: enquiryID.String()}}
	handlers.GetEnquiry(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

}

func TestUpdateEnquiryStatus(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	enquiryID := uuid.New()
	courseID := uuid.New()

	// Enquiry lookup
	mock.ExpectQuery(`SELECT .* FROM "enquiries" WHERE id = \$1 ORDER BY "enquiries"\."id" LIMIT \$2`).
		WithArgs(enquiryID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "interested_course_id", "status", "notes", "created_at"}).
			AddRow(enquiryID, "Prasad Patil", "prasad@example.com", "9876543210", courseID, "new", "Interested", time.Now()))

	// Update status - use permissive pattern for GORM update
	mock.ExpectExec(`UPDATE.*enquiries.*status`).
		WithArgs("contacted", enquiryID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := []byte(`{"status":"contacted"}`)
	req := httptest.NewRequest(http.MethodPatch, "/lms/admissions/enquiries/"+enquiryID.String()+"/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: enquiryID.String()}}
	handlers.UpdateEnquiryStatus(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

}

func TestCreateApplication(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	enquiryID := uuid.New()
	courseID := uuid.New()

	// Enquiry lookup
	mock.ExpectQuery(`SELECT .* FROM "enquiries" WHERE id = \$1 ORDER BY "enquiries"\."id" LIMIT \$2`).
		WithArgs(enquiryID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "interested_course_id", "status", "notes", "created_at"}).
			AddRow(enquiryID, "Prasad Patil", "prasad@example.com", "9876543210", courseID, "new", "Interested", time.Now()))

	// Course lookup
	mock.ExpectQuery(`SELECT .* FROM "courses" WHERE id = \$1 ORDER BY "courses"\."id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "level", "thumbnail_url", "status", "total_seats", "booked_seats", "start_date", "end_date", "meet_link", "created_by", "academic_session_id", "created_at", "updated_at"}).
			AddRow(courseID, "Backend Development", "", "", "", "published", 0, 0, nil, nil, "", uuid.Nil, uuid.Nil, nil, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO .*"applications"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	payload := []byte(`{
		"enquiry_id":"` + enquiryID.String() + `",
		"applied_course_id":"` + courseID.String() + `"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/lms/admissions/applications", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	handlers.CreateApplication(c)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListApplications(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	appID := uuid.New()
	enquiryID := uuid.New()
	courseID := uuid.New()

	// Count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "applications"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// List query
	mock.ExpectQuery(`SELECT .* FROM "applications" LIMIT \$1`).
		WithArgs(20).
		WillReturnRows(sqlmock.NewRows([]string{"id", "enquiry_id", "applied_course_id", "application_status", "remarks", "submitted_at"}).
			AddRow(appID, enquiryID, courseID, "pending", "", time.Now()))

	req := httptest.NewRequest(http.MethodGet, "/lms/admissions/applications", nil)
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	handlers.ListApplications(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetApplication(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	appID := uuid.New()
	enquiryID := uuid.New()
	courseID := uuid.New()

	// Application lookup
	mock.ExpectQuery(`SELECT .* FROM "applications" WHERE id = \$1 ORDER BY "applications"\."id" LIMIT \$2`).
		WithArgs(appID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "enquiry_id", "applied_course_id", "application_status", "remarks", "submitted_at"}).
			AddRow(appID, enquiryID, courseID, "pending", "", time.Now()))

	req := httptest.NewRequest(http.MethodGet, "/lms/admissions/applications/"+appID.String(), nil)
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: appID.String()}}
	handlers.GetApplication(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestApproveApplication(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	appID := uuid.New()
	enquiryID := uuid.New()
	courseID := uuid.New()
	userID := uuid.New()

	// Begin transaction
	mock.ExpectBegin()

	// Application lookup
	mock.ExpectQuery(`SELECT .* FROM "applications" WHERE id = \$1 ORDER BY "applications"\."id" LIMIT \$2`).
		WithArgs(appID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "enquiry_id", "applied_course_id", "application_status", "remarks", "submitted_at"}).
			AddRow(appID, enquiryID, courseID, "pending", "", time.Now()))

	// Enquiry lookup
	mock.ExpectQuery(`SELECT .* FROM "enquiries" WHERE id = \$1 ORDER BY "enquiries"\."id" LIMIT \$2`).
		WithArgs(enquiryID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "email", "phone", "interested_course_id", "status", "notes", "created_at"}).
			AddRow(enquiryID, "Prasad Naik", "prasad@example.com", "9876543210", courseID, "new", "Interested", time.Now()))

	// Student role lookup
	mock.ExpectQuery(`SELECT .* FROM "roles" WHERE name = \$1 ORDER BY "roles"\."id" LIMIT \$2`).
		WithArgs("student", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(userID, "student", "Student role"))

	// Create user
	mock.ExpectExec(`INSERT INTO .*"users"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create student
	mock.ExpectExec(`INSERT INTO .*"students"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Update application status
	mock.ExpectExec(`UPDATE "applications" SET "application_status"=\$1 WHERE "id" = \$2`).
		WithArgs("approved", appID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Update enquiry status
	// Note: Using a more permissive pattern for GORM's Model().Update() SQL generation
	mock.ExpectExec(`UPDATE.*"enquiries".*SET.*"status"`).
		WithArgs("approved", enquiryID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	payload := []byte(`{"temporary_password":"SecurePass123!"}`)
	req := httptest.NewRequest(http.MethodPatch, "/lms/admissions/applications/"+appID.String()+"/approve", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: appID.String()}}
	handlers.ApproveApplication(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if response["application"] == nil {
		t.Fatalf("expected application in response")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRejectApplication(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	appID := uuid.New()
	enquiryID := uuid.New()
	courseID := uuid.New()

	// Application lookup
	mock.ExpectQuery(`SELECT .* FROM "applications" WHERE id = \$1 ORDER BY "applications"\."id" LIMIT \$2`).
		WithArgs(appID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "enquiry_id", "applied_course_id", "application_status", "remarks", "submitted_at"}).
			AddRow(appID, enquiryID, courseID, "pending", "", time.Now()))

	// Update application - use permissive pattern
	mock.ExpectExec(`UPDATE.*applications.*`).
		WithArgs("rejected", "Does not meet requirements", appID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := []byte(`{"remarks":"Does not meet requirements"}`)
	req := httptest.NewRequest(http.MethodPatch, "/lms/admissions/applications/"+appID.String()+"/reject", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: appID.String()}}
	handlers.RejectApplication(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
