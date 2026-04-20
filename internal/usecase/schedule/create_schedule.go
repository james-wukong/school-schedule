// Package schedule contains the use case for creating a schedule.
// It orchestrates the business logic
// and interacts with the repository to persist the new schedule entity.
package schedule

import (
	"context"
	"errors"
	"fmt"

	"github.com/james-wukong/school-schedule/internal/domain/requirement"
	"github.com/james-wukong/school-schedule/internal/domain/room"
	"github.com/james-wukong/school-schedule/internal/domain/schedule"
	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
	"github.com/james-wukong/school-schedule/internal/domain/scheduler/solver"
	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
	"github.com/james-wukong/school-schedule/internal/interface/http/dto"
)

type CreateScheduleUseCase struct {
	reqRepo  requirement.Repository
	roomRepo room.Repository
	tsRepo   timeslot.Repository
	schdRepo schedule.Repository
}

func NewCreateScheduleUseCase(
	reqRepo requirement.Repository,
	roomRepo room.Repository,
	tsRepo timeslot.Repository,
	schdRepo schedule.Repository,
) *CreateScheduleUseCase {
	return &CreateScheduleUseCase{
		reqRepo:  reqRepo,
		roomRepo: roomRepo,
		tsRepo:   tsRepo,
		schdRepo: schdRepo,
	}
}

func (uc *CreateScheduleUseCase) Execute(
	ctx context.Context,
	input dto.CreateScheduleRequest,
) (dto.CreateScheduleResponse, error) {
	rng := solver.NewFastRNG()

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Printf("║     School Class Scheduler v%.1f      ║\n", 1.0)
	fmt.Println("╚══════════════════════════════════════╝")

	// Initialize requirements, rooms and teachers
	// requirements, rooms, _ := solver.Init()
	// 1. Get requirements for a specific verion of a semester
	reqs, err := uc.reqRepo.GetByVersion(
		ctx, input.SchoolID, input.SemesterID, input.Version.InexactFloat64(),
	)
	if err != nil {
		return dto.CreateScheduleResponse{}, err
	}
	if len(reqs) == 0 {
		return dto.CreateScheduleResponse{},
			fmt.Errorf("no requirements found for school: %d, semester: %d, and version, %f",
				input.SchoolID, input.SemesterID, input.Version,
			)
	}

	// 2. Get available rooms
	var schoolRooms []*room.Rooms
	if !input.ExcludeRooms {
		schoolRooms, err = uc.roomRepo.GetBySchoolID(ctx, input.SchoolID)
		if err != nil {
			return dto.CreateScheduleResponse{}, err
		}
	}

	// 3. Get Available Hours
	timeslots, err := uc.tsRepo.GetBySemesterID(ctx, input.SemesterID)
	if err != nil {
		return dto.CreateScheduleResponse{}, err
	}
	tsTimeslots := timeslot.ToTimeslotMap(timeslots)

	// 4. Map to solver models
	var requirements []*model.Requirement
	for _, row := range reqs {
		r, err := row.ToSolverModel()
		if err != nil {
			return dto.CreateScheduleResponse{}, err
		}
		if row.Teacher == nil {
			return dto.CreateScheduleResponse{}, errors.New("no teacher found")
		}
		if len(row.Teacher.Timeslots) == 0 {
			r.Teacher.AvailableTimes = tsTimeslots
		} else {
			r.Teacher.AvailableTimes = timeslot.ToTimeslotMap(row.Teacher.Timeslots)
		}
		requirements = append(requirements, r)
	}
	var rooms []*model.Room
	if !input.ExcludeRooms {
		for _, row := range schoolRooms {
			r := row.ToSolverModel()
			r.AvailableTimes = tsTimeslots
			rooms = append(rooms, r)
		}
	}

	totalSessions := 0
	for _, r := range requirements {
		totalSessions += r.SessionsPerWeek
	}

	fmt.Println("\nSchool overview:")
	fmt.Printf("\nSchool id: %d, Semester id: %d, version: %.2f",
		input.SchoolID, input.SemesterID, input.Version.InexactFloat64())
	fmt.Printf("  Requirements  : %d\n", len(requirements))
	fmt.Printf("  Total sessions: %d\n", totalSessions)
	if !input.ExcludeRooms {
		fmt.Printf("  Rooms         : %d\n", len(rooms))
	}

	// ── Phase 1: Feasibility ──────────────────
	issues := solver.FeasibilityCheck(requirements, rooms,
		solver.ToTimeslots(timeslots),
		input.ExcludeRooms,
	)
	if len(issues) > 0 {
		fmt.Println("\nStopping: feasibility issues found.")
		return dto.CreateScheduleResponse{}, fmt.Errorf("feasibility issues found")
	}

	// ── Phase 2: Greedy Construction ──────────
	fmt.Println("\n=== Phase 2: Greedy Construction ===")
	initial := solver.GreedyConstruct(rng, requirements, rooms,
		solver.ToTimeslots(timeslots),
		input.ExcludeRooms,
	)
	fmt.Printf("  Placed %d / %d sessions\n", len(initial), totalSessions)
	fmt.Printf("  Hard violations : %d\n", solver.HardViolations(initial, input.ExcludeRooms))
	fmt.Printf("  Soft penalty    : %.1f\n", solver.SoftViolations(initial))

	// ── Phase 3: Simulated Annealing ──────────
	optimised := solver.SimulatedAnnealing(
		rng,
		initial,
		rooms,
		solver.ToTimeslots(timeslots),
		input.ExcludeRooms,
		850.0,   // initial temperature
		0.997,   // cooling rate
		100_000, // iterations
	)

	// Generate a new schedule version number and save optimized assignments to database
	var schedules []*schedule.Schedules
	schdVersion := uc.schdRepo.CreateVersionNumber(ctx, input.SemesterID)
	for _, a := range optimised {
		var s schedule.Schedules
		s.RequirementID = int64(a.Requirement.ID)
		s.SchoolID = input.SchoolID
		s.SemesterID = input.SemesterID
		if !input.ExcludeRooms {
			*s.RoomID = int64(a.Room.ID)
		}
		s.TimeslotID = int64(a.Slot.ID)
		s.Status = schedule.StatusDraft
		s.Version = schdVersion
		schedules = append(schedules, &s)
	}
	if err := uc.schdRepo.CreateInBatches(ctx, schedules); err != nil {
		return dto.CreateScheduleResponse{}, err
	}

	// // Save to csv file
	// classFile, err := os.Create("class_report.csv")
	// defer classFile.Close()
	// // Write the UTF-8 BOM bytes first
	// classFile.Write([]byte{0xEF, 0xBB, 0xBF})
	// if err != nil {
	// 	return err
	// }
	// teacherFile, err := os.Create("teacher_report.csv")
	// defer teacherFile.Close()
	// teacherFile.Write([]byte{0xEF, 0xBB, 0xBF})
	// if err != nil {
	// 	return err
	// }
	// reportService := utils.NewClassReportService(uc.reptRepo)
	// teacherService := utils.NewTeacherReportService(uc.reptRepo)
	// if err := reportService.ExportToCSV(
	// 	ctx, classFile, input.SemesterID, input.Version.InexactFloat64(),
	// ); err != nil {
	// 	return err
	// }
	// if err := teacherService.ExportToCSV(
	// 	ctx, teacherFile, input.SemesterID, input.Version.InexactFloat64(),
	// ); err != nil {
	// 	return err
	// }

	// ── Phase 4: Output ───────────────────────
	solver.PrintTimetable(optimised, solver.SampleHeader(tsTimeslots))
	solver.PrintTeacherSchedules(optimised, input.ExcludeRooms)
	solver.PrintSummary(optimised, totalSessions, input.ExcludeRooms)

	return dto.CreateScheduleResponse{
		SchoolID:        input.SchoolID,
		SemesterID:      input.SemesterID,
		ScheduleVersion: schdVersion,
	}, nil
}
