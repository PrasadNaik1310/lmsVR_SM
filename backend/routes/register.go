package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/lms")
	{

		registerAuthRoutes(api) // dont move api.use() line above this line, it will break login page
		api.Use(middleware.AuthMiddleWare())
		registerCompanyRoutes(api)
		RegisterCompanyUserRoutes(api)
		RegisterAdmissionRoutes(api)
		RegisterCourseRoutes(api)
		RegisterCourseModuleRoutes(api)
		RegisterCourseLessonRoutes(api)
		RegisterCourseScheduleRoutes(api)
		RegisterTeacherRoutes(api)
	}
}

func registerCompanyRoutes(api *gin.RouterGroup) {
	company := api.Group("/company")
	company.Use(middleware.AuthMiddleWare())
	{
		registerCompanySessionRoutes(company)
		registerCompanyBatchRoutes(company)

	}
}
