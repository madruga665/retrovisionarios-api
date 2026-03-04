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

func (r *EventRepository) GetAll(ctx context.Context, year int) ([]models.Event, error) {
	query := "SELECT id, date, name, location, flyer FROM events"
	args := []interface{}{}

	if year > 0 {
		query += " WHERE EXTRACT(YEAR FROM date) = $1"
		args = append(args, year)
	}

	query += " ORDER BY date ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.Event, 0, 10)

	for rows.Next() {
		var e models.Event

		err := rows.Scan(&e.ID, &e.Date, &e.Name, &e.Location, &e.Flyer)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	query := "INSERT INTO events (date, name, location, flyer) VALUES ($1, $2, $3, $4) RETURNING id"
	return r.db.QueryRow(ctx, query, event.Date, event.Name, event.Location, event.Flyer).Scan(&event.ID)
}
