package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCourseScheduleRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		courses.POST("/:id/schedules", middleware.RequirePermission("course_schedule.create"), handlers.CreateSchedule)
		courses.GET("/:id/schedules", middleware.RequirePermission("course_schedule.read"), handlers.ListSchedules)
	}

	schedules := api.Group("/schedules")
	{
		schedules.GET("/:id", middleware.RequirePermission("course_schedule.read"), handlers.GetSchedule)
		schedules.PUT("/:id", middleware.RequirePermission("course_schedule.update"), handlers.UpdateSchedule)
		schedules.DELETE("/:id", middleware.RequirePermission("course_schedule.delete"), handlers.DeleteSchedule)
		schedules.POST("/:id/log", middleware.RequirePermission("course_log.create"), handlers.CreateSessionLog)
		schedules.GET("/:id/log", middleware.RequirePermission("course_log.read"), handlers.GetSessionLog)
	}

	courseLogs := api.Group("/course-logs")
	{
		courseLogs.PUT("/:id", middleware.RequirePermission("course_log.update"), handlers.UpdateSessionLog)
	}
}
