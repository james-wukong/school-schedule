// Package school defines the school entity and related value objects.
// It represents how data looks in the database or business rules.
package school

import "time"

// Schools represents the schools table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Schools struct {
	// id is GENERATED ALWAYS. We use <-:false to prevent GORM from
	// including it in INSERT or UPDATE statements.
	ID int64 `gorm:"column:id;primaryKey;<-:false" json:"id"`

	Name            string `gorm:"column:name;not null;unique" json:"name"`
	Code            string `gorm:"column:code;not null;unique" json:"code"`
	Address         string `gorm:"column:address" json:"address"`
	City            string `gorm:"column:city" json:"city"`
	State           string `gorm:"column:state" json:"state"`
	PostalCode      string `gorm:"column:postal_code" json:"postal_code"`
	Country         string `gorm:"column:country" json:"country"`
	Phone           string `gorm:"column:phone" json:"phone"`
	Email           string `gorm:"column:email" json:"email"`
	Website         string `gorm:"column:website" json:"website"`
	EstablishedYear int    `gorm:"column:established_year" json:"established_year"`
	IsActive        bool   `gorm:"column:is_active;default:true" json:"is_active"`

	// CreatedAt is set by the database DEFAULT. We allow read, but restrict write to 'create' only.
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// UpdatedAt is handled automatically by GORM's autoUpdateTime feature.
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

type SchoolFilterEntity struct {
	Email    *string
	Name     *string
	Code     *string
	IsActive *bool
	Page     int
	Limit    int
}

func NewSchool(name, code string, isActive bool) *Schools {
	return &Schools{
		Name:     name,
		Code:     code,
		IsActive: isActive,
	}
}
