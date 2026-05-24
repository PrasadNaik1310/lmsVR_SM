package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/PrasadNaik1310/LMSVR_SM/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthMiddlewareSetsUUIDUserID(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.AuthMiddleWare())
	r.GET("/ping", func(c *gin.Context) {
		value, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "missing user_id"})
			return
		}
		userID, ok := value.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "wrong type"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID.String()})
	})

	userID := uuid.New()
	req, rr := testutil.NewAuthenticatedRequest(t, http.MethodGet, "/ping", userID)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.AuthMiddleWare())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
