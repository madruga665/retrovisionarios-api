package repositories

import (
	"context"
	"retrovisionarios-api/internal/app/v1/events/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetAll() ([]models.Event, error) {
	query := "SELECT id, date, name, flyer FROM events ORDER BY date ASC"
	rows, err := r.db.Query(context.Background(), query)
	var events []models.Event

	if err != nil {
		return events, err
	}

	defer rows.Close()

	for rows.Next() {
		var e models.Event

		err := rows.Scan(&e.ID, &e.Date, &e.Name, &e.Flyer)
		if err != nil {
			return events, err
		}

		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
