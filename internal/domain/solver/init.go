package solver

import (
	"github.com/james-wukong/school-schedule/internal/domain/model"
)

func Init() ([]*model.Requirement, []*model.Room, []*model.Teacher) {
	days := []model.DayOfWeek{model.Monday, model.Tuesday, model.Wednesday,
		model.Thursday, model.Friday}
	slots := []string{"09:00", "10:00", "11:00", "13:00", "14:00", "15:00"}

	// ── Subjects ──
	math := model.NewSubject(model.Subject{ID: 101, Name: "Math", RequiresLab: false})
	english := model.NewSubject(model.Subject{ID: 102, Name: "English", RequiresLab: false})
	science := model.NewSubject(model.Subject{ID: 103, Name: "Science", RequiresLab: true})
	history := model.NewSubject(model.Subject{ID: 104, Name: "History", RequiresLab: false})
	pe := model.NewSubject(model.Subject{ID: 110, Name: "PE", RequiresLab: false})

	// ── Teachers ──
	alice := model.NewTeacher(model.Teacher{ID: 1001, TenantID: 1, Name: "Alice",
		Subjects:       []model.SubjectID{101},
		AvailableTimes: model.AvailableTimeSlots(days, slots),
	})
	bob := model.NewTeacher(model.Teacher{ID: 1002, TenantID: 1, Name: "Bob",
		Subjects:       []model.SubjectID{102},
		AvailableTimes: model.AvailableTimeSlots(days, slots),
	})
	carol := model.NewTeacher(model.Teacher{ID: 1003, TenantID: 1, Name: "Carol",
		Subjects:       []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots(days, slots),
	})
	dave := model.NewTeacher(model.Teacher{ID: 1004, TenantID: 1, Name: "Dave",
		Subjects:       []model.SubjectID{104},
		AvailableTimes: model.AvailableTimeSlots(days, slots),
	})
	eve := model.NewTeacher(model.Teacher{ID: 1005, TenantID: 1, Name: "Eve",
		Subjects:       []model.SubjectID{105},
		AvailableTimes: model.AvailableTimeSlots(days, slots),
	})
	frank := model.NewTeacher(model.Teacher{ID: 1006, TenantID: 1, Name: "Frank",
		Subjects:       []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots(model.RemoveElement(days, model.Friday), slots),
	})

	teachers := []*model.Teacher{alice, bob, carol, dave, eve, frank}

	// ── Classes ──
	cls10A := model.NewSchoolClass(model.SchoolClass{ID: 101, TenantID: 1, StudentCount: 28, Grade: "1", Class: "01"})
	cls10B := model.NewSchoolClass(model.SchoolClass{ID: 102, TenantID: 1, StudentCount: 30, Grade: "1", Class: "02"})
	cls11A := model.NewSchoolClass(model.SchoolClass{ID: 103, TenantID: 1, StudentCount: 25, Grade: "1", Class: "03"})
	cls11B := model.NewSchoolClass(model.SchoolClass{ID: 104, TenantID: 1, StudentCount: 27, Grade: "1", Class: "04"})

	// ── Rooms ──
	rooms := []*model.Room{
		{ID: 101, TenantID: 1, Name: "Room-101", Capacity: 35,
			AvailableTimes: model.AvailableTimeSlots(days, slots),
		},
		{ID: 102, TenantID: 1, Name: "Room-102", Capacity: 35,
			AvailableTimes: model.AvailableTimeSlots(days, slots),
		},
		{ID: 103, TenantID: 1, Name: "Room-103", Capacity: 35,
			AvailableTimes: model.AvailableTimeSlots(days, slots),
		},
		{ID: 201, TenantID: 1, Name: "Room-201", Capacity: 32, Type: model.Lab,
			AvailableTimes: model.AvailableTimeSlots(days, slots),
		},
		{ID: 202, TenantID: 1, Name: "Room-202", Capacity: 32, Type: model.Lab,
			AvailableTimes: model.AvailableTimeSlots(days, slots),
		},
		{ID: 301, TenantID: 1, Name: "Room-301", Capacity: 60, Type: model.GYM,
			AvailableTimes: model.AvailableTimeSlots(days, slots),
		},
	}

	requirements := []*model.Requirement{
		{ID: 10001, SchoolClass: cls10A, Subject: math, Teacher: alice, SessionsPerWeek: 4},
		{ID: 10002, SchoolClass: cls10A, Subject: english, Teacher: bob, SessionsPerWeek: 4},
		{ID: 10003, SchoolClass: cls10A, Subject: science, Teacher: carol, SessionsPerWeek: 3},
		{ID: 10004, SchoolClass: cls10A, Subject: history, Teacher: dave, SessionsPerWeek: 2},
		{ID: 10005, SchoolClass: cls10A, Subject: pe, Teacher: eve, SessionsPerWeek: 2},

		{ID: 10006, SchoolClass: cls10B, Subject: math, Teacher: frank, SessionsPerWeek: 4},
		{ID: 10007, SchoolClass: cls10B, Subject: english, Teacher: bob, SessionsPerWeek: 4},
		{ID: 10008, SchoolClass: cls10B, Subject: science, Teacher: carol, SessionsPerWeek: 3},
		{ID: 10009, SchoolClass: cls10B, Subject: history, Teacher: dave, SessionsPerWeek: 2},
		{ID: 10010, SchoolClass: cls10B, Subject: pe, Teacher: eve, SessionsPerWeek: 2},

		{ID: 10011, SchoolClass: cls11A, Subject: math, Teacher: alice, SessionsPerWeek: 5},
		{ID: 10012, SchoolClass: cls11A, Subject: english, Teacher: bob, SessionsPerWeek: 4},
		{ID: 10013, SchoolClass: cls11A, Subject: science, Teacher: frank, SessionsPerWeek: 3},
		{ID: 10014, SchoolClass: cls11A, Subject: history, Teacher: dave, SessionsPerWeek: 2},
		{ID: 10015, SchoolClass: cls11A, Subject: pe, Teacher: eve, SessionsPerWeek: 2},

		{ID: 10016, SchoolClass: cls11B, Subject: math, Teacher: frank, SessionsPerWeek: 4},
		{ID: 10017, SchoolClass: cls11B, Subject: english, Teacher: bob, SessionsPerWeek: 4},
		{ID: 10018, SchoolClass: cls11B, Subject: science, Teacher: carol, SessionsPerWeek: 3},
		{ID: 10019, SchoolClass: cls11B, Subject: history, Teacher: dave, SessionsPerWeek: 2},
		{ID: 10020, SchoolClass: cls11B, Subject: pe, Teacher: eve, SessionsPerWeek: 2},
	}

	return requirements, rooms, teachers
}
