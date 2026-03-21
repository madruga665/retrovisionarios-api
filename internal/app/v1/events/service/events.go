package service

import (
	"context"
	"retrovisionarios-api/internal/app/v1/events/models"
)

type EventRepository interface {
	GetAll(ctx context.Context, year int, showDeleted bool) ([]models.Event, error)
	GetByID(ctx context.Context, id int) (*models.Event, error)
	Create(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.UpdateEventRequest) error
	Delete(ctx context.Context, id int) error
}

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAll(ctx context.Context, year int, showDeleted bool) ([]models.Event, error) {
	eventList, err := s.repo.GetAll(ctx, year, showDeleted)

	if err != nil {
		return []models.Event{}, err
	}

	return eventList, nil
}

func (s *EventService) GetByID(ctx context.Context, id int) (*models.Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EventService) Update(ctx context.Context, event *models.UpdateEventRequest) error {
	// Verifica se o evento existe antes de tentar atualizar
	if _, err := s.repo.GetByID(ctx, event.ID); err != nil {
		return err
	}
	return s.repo.Update(ctx, event)
}

func (s *EventService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *EventService) Create(ctx context.Context, event *models.Event) error {
	return s.repo.Create(ctx, event)
}
