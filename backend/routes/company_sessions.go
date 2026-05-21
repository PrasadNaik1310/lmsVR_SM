package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func registerCompanySessionRoutes(company *gin.RouterGroup) {
	sessions := company.Group("/sessions")
	{
		sessions.PUT(":session_id/courses/:course_id/assign", middleware.RequirePermission("company.session.assign"), handlers.AssignCourseToSession)
		sessions.GET(":session_id/courses", middleware.RequirePermission("company.session.read"), handlers.ListCoursesBySession)
	}
}
