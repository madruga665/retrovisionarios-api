package repository

import (
	"context"
	"fmt"
	"retrovisionarios-api/internal/app/v1/events/models"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetAll(ctx context.Context, year int, showDeleted bool) ([]models.Event, error) {
	query := "SELECT id, date, name, location, flyer, deleted FROM events"
	args := []interface{}{}

	if !showDeleted {
		query += " WHERE deleted = FALSE"
		if year > 0 {
			query += " AND EXTRACT(YEAR FROM date) = $1"
			args = append(args, year)
		}
	} else if year > 0 {
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

		err := rows.Scan(&e.ID, &e.Date, &e.Name, &e.Location, &e.Flyer, &e.Deleted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return events, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id int) (*models.Event, error) {
	query := "SELECT id, date, name, location, flyer, deleted FROM events WHERE id = $1"
	var e models.Event

	err := r.db.QueryRow(ctx, query, id).Scan(&e.ID, &e.Date, &e.Name, &e.Location, &e.Flyer, &e.Deleted)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *EventRepository) Update(ctx context.Context, event *models.UpdateEventRequest) error {
	query := "UPDATE events SET "
	var args []interface{}
	argCount := 1

	if event.Name != nil {
		query += fmt.Sprintf("name = $%d, ", argCount)
		args = append(args, *event.Name)
		argCount++
	}

	if event.Date != nil {
		query += fmt.Sprintf("date = $%d, ", argCount)
		args = append(args, *event.Date)
		argCount++
	}

	if event.Location != nil {
		query += fmt.Sprintf("location = $%d, ", argCount)
		args = append(args, *event.Location)
		argCount++
	}

	if event.Flyer != nil {
		query += fmt.Sprintf("flyer = $%d, ", argCount)
		args = append(args, *event.Flyer)
		argCount++
	}

	if event.Deleted != nil {
		query += fmt.Sprintf("deleted = $%d, ", argCount)
		args = append(args, *event.Deleted)
		argCount++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, event.ID)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *EventRepository) Delete(ctx context.Context, id int) error {
	query := "UPDATE events SET deleted = TRUE WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	query := "INSERT INTO events (date, name, location, flyer) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := r.db.QueryRow(ctx, query, event.Date, event.Name, event.Location, event.Flyer).Scan(&event.ID); err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	return nil
}
