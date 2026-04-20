// Package app initializes the application routers and dependencies
// for the every module.
package application

import (
	infraPostgres "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"github.com/james-wukong/school-schedule/internal/interface/http/handler"
	schedulerUC "github.com/james-wukong/school-schedule/internal/usecase/schedule"
)

func (a *App) initScheduleRouter() *handler.ScheduleHandler {
	// 1. Repository Layer: Infrastructure implementation of Domain interfaces ---
	// This variable satisfies the user.Repository interface
	reqRepo := infraPostgres.NewRequirementRepository(a.Database.DB, a.Log)
	roomRepo := infraPostgres.NewRoomRepository(a.Database.DB, a.Log)
	tsRepo := infraPostgres.NewTimeslotRepository(a.Database.DB, a.Log)
	schdRepo := infraPostgres.NewScheduleRepository(a.Database.DB, a.Log)

	// 2. UseCase Layer: Business Logic ---
	// THIS IS HOW YOU CREATE THE createScheduleUC VARIABLE
	// We pass the repository into the UseCase constructor
	createUC := schedulerUC.NewCreateScheduleUseCase(reqRepo, roomRepo, tsRepo, schdRepo)

	// 3. Handler Layer: HTTP Handlers ---
	return handler.NewScheduleHandler(createUC)
}
