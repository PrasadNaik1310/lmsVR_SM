package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func registerCompanyBatchRoutes(company *gin.RouterGroup) {
	courses := company.Group("/courses")
	{
		// List courses for authenticated user
		courses.GET("", middleware.RequirePermission("company.batch.read"), handlers.ListCoursesForUser)
		courses.GET(":course_id/batches", middleware.RequirePermission("company.batch.read"), handlers.ListBatchesByCourse)
		courses.POST(":course_id/batches", middleware.RequirePermission("company.batch.create"), handlers.CreateBatchForCourse)
		courses.GET(":course_id/batches/:batch_id", middleware.RequirePermission("company.batch.read"), handlers.GetBatchDetails)
	}
}
