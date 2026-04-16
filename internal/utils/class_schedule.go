package utils

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/report"
	"github.com/james-wukong/school-schedule/internal/infrastructure/repository"
	"github.com/james-wukong/school-schedule/internal/types"
)

// ClassReportService handles the export logic
type ClassReportService struct {
	repo report.Repository
}

// NewClassReportService creates a new instance
func NewClassReportService(repo report.Repository) *ClassReportService {
	return &ClassReportService{repo: repo}
}

// ExportToCSV writes the schedule data to the provided writer
func (s *ClassReportService) ExportToCSV(
	ctx context.Context, w io.Writer, semesterID int64, version float64,
) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 1. Fetch sorted data
	reportService := repository.NewCachedReportRepository(s.repo, nil)
	rows, err := reportService.GetWeeklyClassReport(ctx, semesterID, version)
	if err != nil {
		return err
	}
	// 2. Get max week day
	maxDay := reportService.GetMaxDay(ctx, semesterID, version)

	// 3. Write Headers
	dayIndex := map[int]string{
		2: "Monday", 3: "Tuesday", 4: "Wednesday",
		5: "Thursday", 6: "Friday", 7: "Saturday", 8: "Sunday",
	}
	headers := []string{"Class", "Timeslot"}
	for day := range maxDay {
		if day <= 6 {
			headers = append(headers, dayIndex[day+2])
		}
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 4. Process Rows
	var currentClass string
	var currentStartTime types.ClockTime
	displayRow := make([]string, len(headers))
	for i := range rows {
		if !time.Time(currentStartTime).IsZero() && currentStartTime != rows[i].StartTime {
			// Insert row
			if err := writer.Write(displayRow); err != nil {
				return err
			}
			displayRow = make([]string, len(headers))
		}
		// Insert break line on class change
		if currentClass != "" && currentClass != rows[i].ClassName {
			if err := writer.Write([]string{}); err != nil {
				return err
			}
		}
		// fill display slice
		if displayRow[0] == "" {
			displayRow[0] = fmt.Sprintf("%d (%s)", rows[i].Grade, rows[i].ClassName)
		}
		if displayRow[1] == "" {
			displayRow[1] = time.Time(rows[i].StartTime).Format("15:04")
		}
		displayRow[int(rows[i].DayOfWeek)+1] = fmt.Sprintf("%s-%s",
			rows[i].TeacherName,
			rows[i].SubjectName,
		)

		currentClass = rows[i].ClassName
		currentStartTime = rows[i].StartTime
	}

	return nil
}
