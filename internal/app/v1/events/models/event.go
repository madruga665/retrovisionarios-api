package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date custom type to handle YYYY-MM-DD format
type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(d).Format("2006-01-02"))), nil
}

func (d Date) IsZero() bool {
	return time.Time(d).IsZero()
}

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan Date: %v", value)
	}
	*d = Date(t)
	return nil
}

// Value implements the driver.Valuer interface.
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return time.Time(d).Format("2006-01-02"), nil
}

// Event representa um show ou apresentação da banda.

type Event struct {
	ID       int     `json:"id"`
	Date     Date    `json:"date"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Flyer    *string `json:"flyer"`
}
