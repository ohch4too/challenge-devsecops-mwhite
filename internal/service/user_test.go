package service

import (
	"challenge/internal/domain"
	"errors"
	"testing"
)

type mockUserRepository struct {
	addFunc    func(u *domain.User) error
	getFunc    func(id string) (*domain.User, error)
	listFunc   func() ([]domain.User, error)
	deleteFunc func(id string) error
}

func (m *mockUserRepository) Add(u *domain.User) error {
	if m.addFunc != nil {
		return m.addFunc(u)
	}
	return nil
}

func (m *mockUserRepository) Get(id string) (*domain.User, error) {
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return &domain.User{}, nil
}

func (m *mockUserRepository) List() ([]domain.User, error) {
	if m.listFunc != nil {
		return m.listFunc()
	}
	return []domain.User{}, nil
}

func (m *mockUserRepository) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func TestUserService_AddUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *domain.User
		mockAdd func(u *domain.User) error
		wantErr bool
	}{
		{
			name: "successful add",
			user: &domain.User{
				Firstname: "John",
				Lastname:  "Doe",
				Login:     "johndoe",
				Password:  "password123",
			},
			mockAdd: func(u *domain.User) error { return nil },
			wantErr: false,
		},
		{
			name: "repository error",
			user: &domain.User{
				Firstname: "Jane",
				Lastname:  "Doe",
				Login:     "janedoe",
				Password:  "password123",
			},
			mockAdd: func(u *domain.User) error { return errors.New("database error") },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{addFunc: tt.mockAdd}
			svc := NewUserService(mockRepo)

			if err := svc.AddUser(tt.user); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockGet func(id string) (*domain.User, error)
		wantErr bool
	}{
		{
			name:   "successful get",
			userID: "1",
			mockGet: func(id string) (*domain.User, error) {
				return &domain.User{
					Firstname: "John",
					Lastname:  "Doe",
					Login:     "johndoe",
				}, nil
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: "999",
			mockGet: func(id string) (*domain.User, error) {
				return nil, errors.New("not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{getFunc: tt.mockGet}
			svc := NewUserService(mockRepo)

			user, err := svc.GetUser(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && user == nil {
				t.Error("GetUser() returned nil user")
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	tests := []struct {
		name     string
		mockList func() ([]domain.User, error)
		wantErr  bool
		wantLen  int
	}{
		{
			name: "successful list",
			mockList: func() ([]domain.User, error) {
				return []domain.User{
					{Firstname: "John", Lastname: "Doe", Login: "johndoe"},
					{Firstname: "Jane", Lastname: "Doe", Login: "janedoe"},
				}, nil
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "empty list",
			mockList: func() ([]domain.User, error) {
				return []domain.User{}, nil
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "repository error",
			mockList: func() ([]domain.User, error) {
				return nil, errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{listFunc: tt.mockList}
			svc := NewUserService(mockRepo)

			users, err := svc.ListUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(users) != tt.wantLen {
				t.Errorf("ListUsers() returned %d users, want %d", len(users), tt.wantLen)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		mockDelete func(id string) error
		wantErr    bool
	}{
		{
			name:       "successful delete",
			userID:     "1",
			mockDelete: func(id string) error { return nil },
			wantErr:    false,
		},
		{
			name:       "user not found",
			userID:     "999",
			mockDelete: func(id string) error { return errors.New("not found") },
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{deleteFunc: tt.mockDelete}
			svc := NewUserService(mockRepo)

			if err := svc.DeleteUser(tt.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
