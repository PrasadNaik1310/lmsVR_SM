package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestCreateBatchForCourse(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	courseID := uuid.New()

	// Course lookup with First()
	mock.ExpectQuery(`SELECT .* FROM "courses" WHERE id = \$1 ORDER BY "courses"\."id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "level", "thumbnail_url", "status", "total_seats", "booked_seats", "start_date", "end_date", "meet_link", "created_by", "academic_session_id", "created_at", "updated_at"}).
			AddRow(courseID, "Biology", "", "", "", "published", 0, 0, nil, nil, "", uuid.Nil, uuid.Nil, nil, nil))

	// Batch uniqueness check with First() - expects no rows (error)
	mock.ExpectQuery(`SELECT .* FROM "batches" WHERE course_id = \$1 ORDER BY "batches"\."id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectBegin()

	// Insert batch - keep the SQL expectation simple so the test stays stable.
	mock.ExpectExec(`INSERT INTO .*"batches"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	payload := []byte(`{"batch_name":"Fall 2026","start_date":"2026-09-01","end_date":"2026-12-15","max_students":30,"status":"active"}`)
	req := httptest.NewRequest(http.MethodPost, "/lms/company/courses/"+courseID.String()+"/batches", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "course_id", Value: courseID.String()}}

	handlers.CreateBatchForCourse(c)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if response["batch"] == nil {
		t.Fatalf("expected batch in response, got %s", rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListBatchesByCourse(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	courseID := uuid.New()
	batchID := uuid.New()
	start, _ := time.Parse("2006-01-02", "2026-09-01")
	end, _ := time.Parse("2006-01-02", "2026-12-15")

	// Course lookup with First()
	mock.ExpectQuery(`SELECT .* FROM "courses" WHERE id = \$1 ORDER BY "courses"\."id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "level", "thumbnail_url", "status", "total_seats", "booked_seats", "start_date", "end_date", "meet_link", "created_by", "academic_session_id", "created_at", "updated_at"}).
			AddRow(courseID, "Biology", "", "", "", "published", 0, 0, nil, nil, "", uuid.Nil, uuid.Nil, nil, nil))

	// Count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "batches" WHERE course_id = \$1`).
		WithArgs(courseID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// List query with LIMIT only (page 1, size 10 means offset 0 so GORM doesn't add OFFSET clause)
	mock.ExpectQuery(`SELECT .* FROM "batches" WHERE course_id = \$1 LIMIT \$2`).
		WithArgs(courseID, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "course_id", "batch_name", "start_date", "end_date", "max_students", "status", "created_at"}).
			AddRow(batchID, courseID, "Fall 2026", start, end, 30, "active", time.Now()))

	req := httptest.NewRequest(http.MethodGet, "/lms/company/courses/"+courseID.String()+"/batches?page=1&size=10", nil)
	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{{Key: "course_id", Value: courseID.String()}}

	handlers.ListBatchesByCourse(c)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response struct {
		Batches []models.Batch `json:"batches"`
		Page    int            `json:"page"`
		Size    int            `json:"size"`
		Total   int64          `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(response.Batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(response.Batches))
	}
	if response.Total != 1 {
		t.Fatalf("expected total 1, got %d", response.Total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
