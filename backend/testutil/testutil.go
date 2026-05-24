package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true, WithoutReturning: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm postgres db: %v", err)
	}

	db.DB = gormDB
	return gormDB, mock
}

func NewAuthenticatedRequest(t *testing.T, method, target string, userID uuid.UUID) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret"
		os.Setenv("JWT_SECRET", secret)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign jwt: %v", err)
	}

	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	return req, httptest.NewRecorder()
}
