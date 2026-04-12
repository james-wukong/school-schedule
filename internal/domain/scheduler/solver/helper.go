package solver

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

const (
	TimeLayout        = "2006-01-02 15:04:05"
	TimeMinutesLayout = "15:04"
)

// NewFastRNG returns a thread-local random number generator
// using the high-performance PCG algorithm.
func NewFastRNG() *rand.Rand {
	// PCG requires two uint64 seeds.
	// We use the current time for the first and a fixed constant for the second.
	seed1 := uint64(time.Now().UnixNano())
	seed2 := uint64(time.Now().Unix())

	pcg := rand.NewPCG(seed1, seed2)
	return rand.New(pcg)
}

// isTeacherAvailable checks if the teacher is available at the given time slot
func isTeacherAvailable(teacher *model.Teacher, slot model.TimeSlot) bool {
	if times, ok := teacher.AvailableTimes[slot.Day]; ok {
		for _, t := range times {
			if t == slot.StartTime {
				return true
			}
		}
	}
	return false
}

// isRoomAvailableInDomain checks if the room is available at the given time slot
// based on the current domain (possible assignments)
func isRoomAvailable(room *model.Room, slot model.TimeSlot) bool {
	if times, ok := room.AvailableTimes[slot.Day]; ok {
		for _, t := range times {
			if t == slot.StartTime {
				return true
			}
		}
	}
	return false
}

// dayDifference returns the number of days between two days of the week
func dayDifference(startDay, endDay model.DayOfWeek) int {
	if endDay < startDay {
		return 10
	}
	return int(endDay) - int(startDay)
}

func minuteDifference(start, end string) int {
	s, err := time.Parse(TimeMinutesLayout, start)
	if err != nil {
		return 0
	}
	e, err := time.Parse(TimeMinutesLayout, end)
	if err != nil {
		return 0
	}
	return int(math.Abs(e.Sub(s).Minutes()))
}

func isTimeAfter(t, benchmark string) (bool, error) {
	t1, err := time.Parse(TimeMinutesLayout, t)
	if err != nil {
		return false, err
	}
	t2, err := time.Parse(TimeMinutesLayout, benchmark)
	if err != nil {
		return false, err
	}
	return t1.After(t2), nil
}

func ToTimeslots(ts map[model.DayOfWeek][]string) []model.TimeSlot {
	var timeslots []model.TimeSlot
	for key, value := range ts {
		for _, slot := range value {
			timeslots = append(timeslots, model.TimeSlot{
				StartTime: slot,
				Day:       key,
			})
		}
	}
	return timeslots
}

func SampleHeader(ts map[model.DayOfWeek][]string) []string {
	return ts[model.Monday]
}
