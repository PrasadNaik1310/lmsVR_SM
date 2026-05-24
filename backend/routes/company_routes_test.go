package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/routes"
	"github.com/PrasadNaik1310/LMSVR_SM/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestBatchListRouteIsWiredAndProtected(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	roleID := uuid.New()
	permissionID := uuid.New()
	courseID := uuid.New()
	batchID := uuid.New()
	start, _ := time.Parse("2006-01-02", "2026-09-01")
	end, _ := time.Parse("2006-01-02", "2026-12-15")

	// User lookup for auth middleware
	mock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "phone", "role_id", "is_active", "last_login_at", "created_at", "updated_at"}).
			AddRow(userID, "Test", "User", "test@example.com", "secret", "", roleID, true, nil, nil, nil))

	// Permission lookup
	mock.ExpectQuery(`SELECT .* FROM "permissions" WHERE name = \$1 ORDER BY "permissions"\."id" LIMIT \$2`).
		WithArgs("company.batch.read", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(permissionID, "company.batch.read", "test permission"))

	// Role-permission link check
	mock.ExpectQuery(`SELECT .* FROM "role_permissions" WHERE role_id = \$1 AND permission_id = \$2 ORDER BY "role_permissions"\."id" LIMIT \$3`).
		WithArgs(roleID, permissionID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "permission_id"}).AddRow(uuid.New(), roleID, permissionID))

	// Course lookup
	mock.ExpectQuery(`SELECT .* FROM "courses" WHERE id = \$1 ORDER BY "courses"\."id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "level", "thumbnail_url", "status", "total_seats", "booked_seats", "start_date", "end_date", "meet_link", "created_by", "academic_session_id", "created_at", "updated_at"}).
			AddRow(courseID, "Biology", "", "", "", "published", 0, 0, nil, nil, "", uuid.Nil, uuid.Nil, nil, nil))

	// Batch count
	mock.ExpectQuery(`SELECT count\(\*\) FROM "batches" WHERE course_id = \$1`).
		WithArgs(courseID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Batch list with pagination (default size=20, page 1 means offset=0 so GORM only adds LIMIT)
	mock.ExpectQuery(`SELECT .* FROM "batches" WHERE course_id = \$1 LIMIT \$2`).
		WithArgs(courseID, 20).
		WillReturnRows(sqlmock.NewRows([]string{"id", "course_id", "batch_name", "start_date", "end_date", "max_students", "status", "created_at"}).
			AddRow(batchID, courseID, "Fall 2026", start, end, 30, "active", time.Now()))

	r := gin.New()
	routes.RegisterRoutes(r)

	req, rr := testutil.NewAuthenticatedRequest(t, http.MethodGet, "/lms/company/courses/"+courseID.String()+"/batches", userID)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response struct {
		Batches []models.Batch `json:"batches"`
		Total   int64          `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(response.Batches) != 1 || response.Batches[0].CourseID != courseID {
		t.Fatalf("unexpected batches response: %+v", response)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCompanyRoutesRejectMissingAuth(t *testing.T) {
	_, _ = testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	routes.RegisterRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/lms/company/courses/"+uuid.New().String()+"/batches", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
