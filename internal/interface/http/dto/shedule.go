package dto

type CreateScheduleRequest struct {
	SemesterID   int64
	SchoolID     int64
	Version      float64
	ExcludeRooms bool
}

type CreateScheduleResponse struct {
}
