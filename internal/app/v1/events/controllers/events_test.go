package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"retrovisionarios-api/internal/app/v1/events/models"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockEventService struct {
	GetAllFunc func(year int) ([]models.Event, error)
	CreateFunc func(event *models.Event) error
}

func (m *mockEventService) GetAll(year int) ([]models.Event, error) {
	return m.GetAllFunc(year)
}

func (m *mockEventService) Create(event *models.Event) error {
	return m.CreateFunc(event)
}

func TestEventController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

		t.Run("Success Happy Path", func(t *testing.T) {
			mockService := &mockEventService{
				CreateFunc: func(event *models.Event) error {
					event.ID = 1
					return nil
				},
			}
			controller := NewEventController(mockService)
	
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.POST("/events", controller.Create)
	
			eventDate := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)
			input := map[string]interface{}{
				"name": "Evento Teste",
				"date": "2026-02-20",
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
			if !strings.Contains(w.Body.String(), "Os campos 'date' e 'name' são obrigatórios") {
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
			if !strings.Contains(w.Body.String(), "Payload malformado ou campos inválidos") {
				t.Errorf("Expected error message about malformed payload, got %s", w.Body.String())
			}
		})
	
		t.Run("Internal Server Error", func(t *testing.T) {
			mockService := &mockEventService{
				CreateFunc: func(event *models.Event) error {
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
