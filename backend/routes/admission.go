package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdmissionRoutes(api *gin.RouterGroup) {
	admissions := api.Group("/admissions")
	admissions.Use(middleware.AuthMiddleWare())
	{
		registerAdmissionEnquiryRoutes(admissions)
		registerAdmissionApplicationRoutes(admissions)
	}
}

func registerAdmissionEnquiryRoutes(group *gin.RouterGroup) {
	enquiries := group.Group("/enquiries")
	{
		enquiries.POST("", middleware.RequirePermission("admission.enquiry.create"), handlers.CreateEnquiry)
		enquiries.GET("", middleware.RequirePermission("admission.enquiry.read"), handlers.ListEnquiries)
		enquiries.GET("/:id", middleware.RequirePermission("admission.enquiry.read"), handlers.GetEnquiry)
		enquiries.PATCH("/:id/status", middleware.RequirePermission("admission.enquiry.update"), handlers.UpdateEnquiryStatus)
	}
}

func registerAdmissionApplicationRoutes(group *gin.RouterGroup) {
	applications := group.Group("/applications")
	{
		applications.POST("", middleware.RequirePermission("admission.application.create"), handlers.CreateApplication)
		applications.GET("", middleware.RequirePermission("admission.application.read"), handlers.ListApplications)
		applications.GET("/:id", middleware.RequirePermission("admission.application.read"), handlers.GetApplication)
		applications.PATCH("/:id/approve", middleware.RequirePermission("admission.application.approve"), handlers.ApproveApplication)
		applications.PATCH("/:id/reject", middleware.RequirePermission("admission.application.reject"), handlers.RejectApplication)
	}
}
