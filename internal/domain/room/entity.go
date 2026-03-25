// Package room defines the rooms entity and related value objects.
// It represents how data looks in the database or business rules.
package room

import (
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/school"
)

type RoomType string

const (
	Regular RoomType = "Regular"
	Lab     RoomType = "Lab"
	Gym     RoomType = "Gym"
)

// Rooms represents the rooms table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Rooms struct {
	// ID uses the Postgres IDENTITY sequence. default:nextval handles the START WITH 1000 logic.
	ID int64 `gorm:"primaryKey;column:id;default:nextval('rooms_id_seq');<-:false" json:"id"`

	// Foreign Key to School
	SchoolID int64           `gorm:"column:school_id;not null;uniqueIndex:idx_rooms_school_code" json:"school_id"`
	School   *school.Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`

	// Room Identity within the School
	Code string `gorm:"column:code;type:varchar(50);not null;uniqueIndex:idx_rooms_school_code" json:"code"`
	Name string `gorm:"column:name;type:varchar(100);not null" json:"name"`

	// Classification and Capacity
	RoomType    RoomType `gorm:"column:room_type;type:varchar(50);index:idx_rooms_type" json:"room_type"`
	Capacity    int      `gorm:"column:capacity;default:40" json:"capacity"`
	FloorNumber *int     `gorm:"column:floor_number" json:"floor_number"` // Pointer to allow for NULL
	Building    string   `gorm:"column:building;type:varchar(100)" json:"building"`

	// Scheduling Constraints
	AvailableDays string `gorm:"column:available_days;type:varchar(50)" json:"available_days"`
	IsActive      bool   `gorm:"column:is_active;not null;default:true;index:idx_rooms_active" json:"is_active"`

	// Audit Timestamps
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func NewRoom(schoolID int64,
	capacity int,
	code, name string,
	isActive bool,
	roomType RoomType,
) *Rooms {
	return &Rooms{
		SchoolID: schoolID,
		Code:     code,
		Name:     name,
		RoomType: roomType,
		Capacity: capacity,
		IsActive: isActive,
	}
}
