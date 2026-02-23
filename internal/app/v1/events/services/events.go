package services

import (
	"context"
	"retrovisionarios-api/internal/app/v1/events/models"
)

type EventRepository interface {
	GetAll(ctx context.Context, year int) ([]models.Event, error)
	Create(ctx context.Context, event *models.Event) error
}

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAll(ctx context.Context, year int) ([]models.Event, error) {
	eventList, err := s.repo.GetAll(ctx, year)

	if err != nil {
		return []models.Event{}, err
	}

	return eventList, nil
}

func (s *EventService) Create(ctx context.Context, event *models.Event) error {
	return s.repo.Create(ctx, event)
}
