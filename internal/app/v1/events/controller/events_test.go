package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"retrovisionarios-api/internal/app/v1/events/models"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type mockEventService struct {
	GetAllFunc  func(ctx context.Context, year int, showDeleted bool, name string) ([]models.Event, error)
	GetByIDFunc func(ctx context.Context, id int) (*models.Event, error)
	CreateFunc  func(ctx context.Context, event *models.Event) error
	UpdateFunc  func(ctx context.Context, event *models.UpdateEventRequest) error
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *mockEventService) GetAll(ctx context.Context, year int, showDeleted bool, name string) ([]models.Event, error) {
	return m.GetAllFunc(ctx, year, showDeleted, name)
}

func (m *mockEventService) GetByID(ctx context.Context, id int) (*models.Event, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *mockEventService) Create(ctx context.Context, event *models.Event) error {
	return m.CreateFunc(ctx, event)
}

func (m *mockEventService) Update(ctx context.Context, event *models.UpdateEventRequest) error {
	return m.UpdateFunc(ctx, event)
}

func (m *mockEventService) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

func TestEventController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Happy Path", func(t *testing.T) {
		mockService := &mockEventService{
			CreateFunc: func(ctx context.Context, event *models.Event) error {
				event.ID = 1
				return nil
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.POST("/events", controller.Create)

		eventDate := time.Date(2026, 2, 20, 21, 30, 0, 0, time.UTC)
		input := map[string]interface{}{
			"name":     "Evento Teste",
			"date":     "2026-02-20 21:30",
			"location": "Local Teste",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var response models.Event
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("Expected ID 1, got %d", response.ID)
		}
		if response.Name != "Evento Teste" {
			t.Errorf("Expected name 'Evento Teste', got '%s'", response.Name)
		}
		if response.Location != "Local Teste" {
			t.Errorf("Expected location 'Local Teste', got '%s'", response.Location)
		}
		if !time.Time(response.Date).Equal(eventDate) {
			t.Errorf("Expected date %v, got %v", eventDate, response.Date)
		}
	})

	t.Run("Validation Failure - Missing Name", func(t *testing.T) {
		mockService := &mockEventService{}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.POST("/events", controller.Create)

		input := map[string]interface{}{
			"date": "2026-02-20",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "os campos 'date' e 'name' são obrigatórios") {
			t.Errorf("Expected error message about required fields, got %s", w.Body.String())
		}
	})

	t.Run("Type Failure - Invalid Date", func(t *testing.T) {
		mockService := &mockEventService{}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.POST("/events", controller.Create)

		input := map[string]interface{}{
			"name": "Evento Teste",
			"date": "20-02-2026", // Wrong format
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "os campos 'date' e 'name' são obrigatórios e devem estar no formato correto") {
			t.Errorf("Expected error message about malformed payload, got %s", w.Body.String())
		}
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockService := &mockEventService{
			CreateFunc: func(ctx context.Context, event *models.Event) error {
				return errors.New("db error")
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.POST("/events", controller.Create)

		input := map[string]interface{}{
			"name": "Evento Teste",
			"date": "2026-02-20",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "Erro interno ao criar evento") {
			t.Errorf("Expected error message about internal error, got %s", w.Body.String())
		}
	})

}

func TestEventController_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - List all active events", func(t *testing.T) {
		mockService := &mockEventService{
			GetAllFunc: func(ctx context.Context, year int, showDeleted bool, name string) ([]models.Event, error) {
				return []models.Event{
					{ID: 1, Name: "Evento 1", Deleted: false},
					{ID: 2, Name: "Evento 2", Deleted: false},
				}, nil
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/events", controller.GetAll)

		req, _ := http.NewRequest(http.MethodGet, "/events", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string][]models.Event
		json.Unmarshal(w.Body.Bytes(), &response)
		if len(response["result"]) != 2 {
			t.Errorf("Expected 2 events, got %d", len(response["result"]))
		}
	})

	t.Run("Success - Filter by year and deleted", func(t *testing.T) {
		mockService := &mockEventService{
			GetAllFunc: func(ctx context.Context, year int, showDeleted bool, name string) ([]models.Event, error) {
				if year == 2026 && showDeleted == true {
					return []models.Event{
						{ID: 1, Name: "Evento Deletado", Deleted: true},
					}, nil
				}
				return []models.Event{}, nil
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/events", controller.GetAll)

		req, _ := http.NewRequest(http.MethodGet, "/events?year=2026&deleted=true", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string][]models.Event
		json.Unmarshal(w.Body.Bytes(), &response)
		if len(response["result"]) != 1 {
			t.Errorf("Expected 1 event, got %d", len(response["result"]))
		}
		if response["result"][0].Deleted != true {
			t.Errorf("Expected deleted event")
		}
	})

	t.Run("Success - Filter by name", func(t *testing.T) {
		mockService := &mockEventService{
			GetAllFunc: func(ctx context.Context, year int, showDeleted bool, name string) ([]models.Event, error) {
				if name == "Evento 1" {
					return []models.Event{
						{ID: 1, Name: "Evento 1"},
					}, nil
				}
				return []models.Event{}, nil
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/events", controller.GetAll)

		req, _ := http.NewRequest(http.MethodGet, "/events?name=Evento 1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string][]models.Event
		json.Unmarshal(w.Body.Bytes(), &response)
		if len(response["result"]) != 1 {
			t.Errorf("Expected 1 event, got %d", len(response["result"]))
		}
		if response["result"][0].Name != "Evento 1" {
			t.Errorf("Expected event name 'Evento 1', got '%s'", response["result"][0].Name)
		}
	})
}

func TestEventController_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Happy Path", func(t *testing.T) {
		mockService := &mockEventService{
			DeleteFunc: func(ctx context.Context, id int) error {
				if id == 1 {
					return nil
				}
				return errors.New("event not found")
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.DELETE("/events/:id", controller.Delete)

		req, _ := http.NewRequest(http.MethodDelete, "/events/1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		controller := NewEventController(&mockEventService{})

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.DELETE("/events/:id", controller.Delete)

		req, _ := http.NewRequest(http.MethodDelete, "/events/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestEventController_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Happy Path", func(t *testing.T) {
		mockService := &mockEventService{
			UpdateFunc: func(ctx context.Context, event *models.UpdateEventRequest) error {
				return nil
			},
			GetByIDFunc: func(ctx context.Context, id int) (*models.Event, error) {
				return &models.Event{ID: id, Name: "Evento Atualizado"}, nil
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.PATCH("/events/:id", controller.Update)

		input := map[string]interface{}{
			"name":     "Evento Atualizado",
			"date":     "2026-03-10 20:00",
			"location": "Novo Local",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPatch, "/events/1", bytes.NewBuffer(body))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response models.Event
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != "Evento Atualizado" {
			t.Errorf("Expected name 'Evento Atualizado', got '%s'", response.Name)
		}
	})

	t.Run("Event Not Found", func(t *testing.T) {
		mockService := &mockEventService{
			UpdateFunc: func(ctx context.Context, event *models.UpdateEventRequest) error {
				return pgx.ErrNoRows
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.PATCH("/events/:id", controller.Update)

		input := map[string]interface{}{
			"name": "Evento Inexistente",
			"date": "2026-03-10 20:00",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPatch, "/events/999", bytes.NewBuffer(body))

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "Evento não encontrado.") {
			t.Errorf("Expected error message about event not found, got %s", w.Body.String())
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		controller := NewEventController(&mockEventService{})

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.PATCH("/events/:id", controller.Update)

		req, _ := http.NewRequest(http.MethodPatch, "/events/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestEventController_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Happy Path", func(t *testing.T) {
		eventDate := time.Date(2026, 3, 25, 20, 0, 0, 0, time.UTC)
		mockService := &mockEventService{
			GetByIDFunc: func(ctx context.Context, id int) (*models.Event, error) {
				if id == 1 {
					return &models.Event{
						ID:       1,
						Name:     "Evento Teste",
						Date:     models.DateTime(eventDate),
						Location: "Local Teste",
					}, nil
				}
				return nil, errors.New("no rows in result set")
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/events/:id", controller.GetByID)

		req, _ := http.NewRequest(http.MethodGet, "/events/1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response models.Event
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("Expected ID 1, got %d", response.ID)
		}
		if response.Name != "Evento Teste" {
			t.Errorf("Expected name 'Evento Teste', got '%s'", response.Name)
		}
	})

	t.Run("Event Not Found", func(t *testing.T) {
		mockService := &mockEventService{
			GetByIDFunc: func(ctx context.Context, id int) (*models.Event, error) {
				return nil, pgx.ErrNoRows
			},
		}
		controller := NewEventController(mockService)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/events/:id", controller.GetByID)

		req, _ := http.NewRequest(http.MethodGet, "/events/999", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "Evento não encontrado.") {
			t.Errorf("Expected error message about event not found, got %s", w.Body.String())
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		controller := NewEventController(&mockEventService{})

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/events/:id", controller.GetByID)

		req, _ := http.NewRequest(http.MethodGet, "/events/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}
