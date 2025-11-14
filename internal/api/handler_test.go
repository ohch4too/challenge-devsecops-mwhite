package api

import (
	"bytes"
	"challenge/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockUserService struct {
	addFunc    func(u *domain.User) error
	getFunc    func(id string) (*domain.User, error)
	listFunc   func() ([]domain.User, error)
	deleteFunc func(id string) error
}

func (m *mockUserService) AddUser(u *domain.User) error {
	if m.addFunc != nil {
		return m.addFunc(u)
	}
	return nil
}

func (m *mockUserService) GetUser(id string) (*domain.User, error) {
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return &domain.User{}, nil
}

func (m *mockUserService) ListUsers() ([]domain.User, error) {
	if m.listFunc != nil {
		return m.listFunc()
	}
	return []domain.User{}, nil
}

func (m *mockUserService) DeleteUser(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockList       func() ([]domain.User, error)
		expectedStatus int
	}{
		{
			name: "successful list",
			mockList: func() ([]domain.User, error) {
				return []domain.User{
					{Firstname: "John", Lastname: "Doe", Login: "johndoe"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error listing",
			mockList: func() ([]domain.User, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{listFunc: tt.mockList}
			handler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/users", nil)

			handler.ListUsers(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListUsers() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_AddUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           interface{}
		mockAdd        func(u *domain.User) error
		expectedStatus int
	}{
		{
			name: "successful add",
			body: map[string]string{
				"firstname": "John",
				"lastname":  "Doe",
				"login":     "johndoe",
				"password":  "password123",
			},
			mockAdd:        func(u *domain.User) error { return nil },
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			body: map[string]string{
				"firstname": "",
				"lastname":  "Doe",
				"login":     "johndoe",
				"password":  "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{addFunc: tt.mockAdd}
			handler := NewUserHandler(mockSvc)

			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.AddUser(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("AddUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockGet        func(id string) (*domain.User, error)
		expectedStatus int
	}{
		{
			name:   "successful get",
			userID: "1",
			mockGet: func(id string) (*domain.User, error) {
				return &domain.User{Firstname: "John", Lastname: "Doe", Login: "johndoe"}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "999",
			mockGet: func(id string) (*domain.User, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{getFunc: tt.mockGet}
			handler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.userID}}

			handler.GetUser(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestUserHandler_DelUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockDelete     func(id string) error
		expectedStatus int
	}{
		{
			name:           "successful delete",
			userID:         "1",
			mockDelete:     func(id string) error { return nil },
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "user not found",
			userID: "999",
			mockDelete: func(id string) error {
				return errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{deleteFunc: tt.mockDelete}
			handler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.userID}}

			handler.DelUser(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("DelUser() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}
