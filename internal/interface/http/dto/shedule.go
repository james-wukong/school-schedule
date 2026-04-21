// Package dto
package dto

import (
	"github.com/shopspring/decimal"
)

type CreateScheduleRequest struct {
	SemesterID   int64           `json:"semester_id"`
	SchoolID     int64           `json:"school_id"`
	Version      decimal.Decimal `json:"version"`
	ExcludeRooms bool            `json:"exclude_rooms"`
}

type CreateScheduleResponse struct {
	SemesterID      int64           `json:"semester_id"`
	SchoolID        int64           `json:"school_id"`
	ScheduleVersion decimal.Decimal `json:"schedule_version"`
}

type ExportAPIResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}
