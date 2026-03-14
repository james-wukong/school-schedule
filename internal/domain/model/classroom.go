package model

// type Classroom struct {
// 	ID       string
// 	Name     string
// 	Capacity int
// 	// AvailableTimes: map of day -> available times
// 	AvailableTimes map[string][]string
// }

type Classroom struct {
	// AvailableTimes: map of day -> available times
	AvailableTimes map[DayOfWeek][]string
	Name           string
	Equipment      []string
	ID             int
	TenantID       int
	Capacity       int
}

func (c Classroom) CanFit(size int) bool {
	return c.Capacity >= size
}
