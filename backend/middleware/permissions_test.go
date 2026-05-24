package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/PrasadNaik1310/LMSVR_SM/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequirePermissionAllowsAuthorizedUser(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	roleID := uuid.New()
	permissionID := uuid.New()
	permissionName := "company.batch.read"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "phone", "role_id", "is_active", "last_login_at", "created_at", "updated_at"}).
			AddRow(userID, "Test", "User", "test@example.com", "secret", "", roleID, true, nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE name = $1 ORDER BY "permissions"."id" LIMIT $2`)).
		WithArgs(permissionName, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(permissionID, permissionName, "test permission"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE role_id = $1 AND permission_id = $2 ORDER BY "role_permissions"."id" LIMIT $3`)).
		WithArgs(roleID, permissionID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "permission_id"}).AddRow(uuid.New(), roleID, permissionID))

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	r.Use(middleware.RequirePermission(permissionName))
	r.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRequirePermissionDeniesUnauthorizedUser(t *testing.T) {
	_, mock := testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	roleID := uuid.New()
	permissionID := uuid.New()
	permissionName := "company.batch.read"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "phone", "role_id", "is_active", "last_login_at", "created_at", "updated_at"}).
			AddRow(userID, "Test", "User", "denied@example.com", "secret", "", roleID, true, nil, nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE name = $1 ORDER BY "permissions"."id" LIMIT $2`)).
		WithArgs(permissionName, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(permissionID, permissionName, "test permission"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE role_id = $1 AND permission_id = $2 ORDER BY "role_permissions"."id" LIMIT $3`)).
		WithArgs(roleID, permissionID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "permission_id"}))

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	r.Use(middleware.RequirePermission(permissionName))
	r.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRequirePermissionRejectsMissingUserID(t *testing.T) {
	_, _ = testutil.SetupMockDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.RequirePermission("company.batch.read"))
	r.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
