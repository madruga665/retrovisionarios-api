package models

import (
	"time"
)

// Event representa um show ou apresentação da banda.
type Event struct {
	ID    int       `json:"id"`
	Date  time.Time `json:"date"`
	Name  string    `json:"name"`
	Flyer string    `json:"flyer"`
}
