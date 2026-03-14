package http

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/james-wukong/school-schedule/internal/application/usecase"
// )

// func RegisterRoutes(r *gin.Engine, uc *usecase.ScheduleUseCase) {
// 	r.GET("/health", healthHandler)
// 	r.GET("/ready", readinessHandler)

// 	v1 := r.Group("/v1")
// 	v1.Use(TenantMiddleware())

// 	schedules := v1.Group("/schedules")
// 	{
// 		schedules.POST("/generate", generateScheduleHandler(uc))
// 		schedules.GET("/:scheduleID", getScheduleHandler(uc))
// 		schedules.GET("/jobs/:jobID", getJobStatusHandler(uc))
// 		schedules.POST("/:scheduleID/publish", publishScheduleHandler(uc))
// 	}
// }
