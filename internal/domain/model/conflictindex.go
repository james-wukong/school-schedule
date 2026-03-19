package model

import (
	"fmt"
)

type ConflictIndex struct {
	TeacherSlot map[string][]*Assignment // "teacherID_day_time"
	ClassSlot   map[string][]*Assignment // "classID_day_time"
	RoomSlot    map[string][]*Assignment // "roomID_day_time"
}

type ConflictKeyID interface {
	// Use tilde (~) to allow underlying types like 'type TeacherID int'
	~int
}

type ConflictMapKey[T ConflictKeyID] struct {
	ID   T
	Day  DayOfWeek
	Slot string
}

func NewConflictIndex() *ConflictIndex {
	return &ConflictIndex{
		TeacherSlot: make(map[string][]*Assignment), // is any teacher in two places at once
		ClassSlot:   make(map[string][]*Assignment), // is any class having two lessons at once?
		RoomSlot:    make(map[string][]*Assignment), // is any room hosting two classes at once?
	}
}

func ConflictKey[T TeacherID | ClassID | RoomID](id T, slot TimeSlot) string {
	return fmt.Sprintf("%d_%d_%s", id, slot.Day, slot.StartTime)
}
