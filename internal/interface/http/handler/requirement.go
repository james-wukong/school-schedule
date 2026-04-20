// Package handler contains HTTP handlers for schedule-related endpoints.
package handler

import (
	"net/http"

	"github.com/james-wukong/school-schedule/internal/interface/http/dto"
	"github.com/james-wukong/school-schedule/internal/interface/http/middleware"
	"github.com/james-wukong/school-schedule/internal/usecase/schedule"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	createScheduleUC *schedule.CreateScheduleUseCase
	// getScheduleUC    *schedule.GetScheduleUseCase
}

func NewScheduleHandler(
	c *schedule.CreateScheduleUseCase,
	// g *schedule.GetScheduleUseCase,
) *ScheduleHandler {
	return &ScheduleHandler{
		createScheduleUC: c,
	}
}

// Register satisfies the RouterRegister interface
func (h *ScheduleHandler) Register(mw *middleware.Manager, v1 *gin.RouterGroup) {
	resGroup := v1.Group("/schedule")
	{
		resGroup.POST("/create", h.Create)
		// resGroup.GET("/:id", h.GetProfile)
	}
}

// Create godoc
//
//	@Summary		Create
//	@Description	create a schedule
//	@Tags			create schedule
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateScheduleRequest	true	"User Credentials"
//	@Success		201		{object}	dto.CreateScheduleResponse
//	@Failure		400		{object}	map[string]any	"{'error':'error message'}"
//	@Failure		500		{object}	map[string]any	"{'error':'error message'}"
//	@Router			/schedule/create [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.createScheduleUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Map domain entity to response DTO
	c.JSON(http.StatusCreated, nil)
}
