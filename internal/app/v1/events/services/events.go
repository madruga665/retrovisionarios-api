package services

import "retrovisionarios-api/internal/app/v1/events/models"

type EventRepository interface {
	GetAll(year int) ([]models.Event, error)
}

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAll(year int) ([]models.Event, error) {
	eventList, err := s.repo.GetAll(year)

	if err != nil {
		return []models.Event{}, err
	}

	return eventList, nil
}
