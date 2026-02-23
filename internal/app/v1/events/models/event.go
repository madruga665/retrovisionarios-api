package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// DateTime custom type to handle YYYY-MM-DD HH:MM format
type DateTime time.Time

func (d *DateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	// Suporta tanto o formato com hora quanto apenas data para compatibilidade
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
	}
	*d = DateTime(t)
	return nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(d).Format("2006-01-02 15:04"))), nil
}

func (d DateTime) IsZero() bool {
	return time.Time(d).IsZero()
}

// Scan implements the sql.Scanner interface.
func (d *DateTime) Scan(value interface{}) error {
	if value == nil {
		*d = DateTime(time.Time{})
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan DateTime: %v", value)
	}
	*d = DateTime(t)
	return nil
}

// Value implements the driver.Valuer interface.
func (d DateTime) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return time.Time(d), nil
}

// Event representa um show ou apresentação da banda.

type Event struct {
	ID       int      `json:"id"`
	Date     DateTime `json:"date"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Flyer    *string  `json:"flyer"`
}
