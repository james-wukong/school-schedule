// Package types define some custom types of time
package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CivilDate time.Time

const dateFormat = "2006-01-02"

// UnmarshalJSON JSON Support
func (c *CivilDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	*c = CivilDate(t)
	return nil
}

// UnmarshalCSV CSV Support (csvutil uses UnmarshalCSV or TextUnmarshaler)
func (c *CivilDate) UnmarshalCSV(data []byte) error {
	return c.UnmarshalText(data)
}

// UnmarshalText Form/Text Support (Used by many Go form decoders)
func (c *CivilDate) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	t, err := time.Parse(dateFormat, string(data))
	if err != nil {
		return err
	}
	*c = CivilDate(t)
	return nil
}

// Value GORM/SQL Support (So you can save it directly to the DB)
func (c CivilDate) Value() (driver.Value, error) {
	if time.Time(c).IsZero() {
		return nil, nil
	}
	return time.Time(c).Format(dateFormat), nil
}

func (c *CivilDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*c = CivilDate(v)
	case string:
		t, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
		*c = CivilDate(t)
	case []byte:
		t, err := time.Parse(dateFormat, string(v))
		if err != nil {
			return err
		}
		*c = CivilDate(t)
	default:
		return fmt.Errorf("cannot scan %T into CivilDate", value)
	}

	return nil
}

type Int64Slice []int64

// UnmarshalCSV converts a comma-separated string "1,2,3" into []int64
func (s *Int64Slice) UnmarshalCSV(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Split by comma
	cleanData := strings.Trim(string(data), "\" ")
	parts := strings.Split(cleanData, ",")
	for _, p := range parts {
		val, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			return err
		}
		*s = append(*s, val)
	}
	return nil
}

type ClockTime time.Time

const clockFormat = "15:04"

// UnmarshalJSON handles the "09:00" -> HourMinute conversion
// When json.Unmarshal encounters a field of type HourMinute,
// it checks if that type has an UnmarshalJSON([]byte) error method
// UnmarshalJSON handles JSON: {"start_time": "15:13"}
func (c *ClockTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	return c.parse(s)
}

// UnmarshalText handles CSV and Form mapping
func (c *ClockTime) UnmarshalText(b []byte) error {
	s := string(b)
	if s == "" {
		return nil
	}
	return c.parse(s)
}

// Internal helper to keep code DRY
func (c *ClockTime) parse(s string) error {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		t, err = time.Parse(clockFormat, s)
		if err != nil {
			return fmt.Errorf("invalid time format (expected HH:MM): %w", err)
		}
	}
	*c = ClockTime(t)
	return nil
}

func (c ClockTime) Value() (driver.Value, error) {
	return time.Time(c).Format(clockFormat), nil
}

func (c *ClockTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	// Most DB drivers return TIME columns as time.Time or string
	switch v := value.(type) {
	case time.Time:
		*c = ClockTime(v)
	case string:
		return c.parse(v)
	case []byte:
		return c.parse(string(v))
	}
	return nil
}
