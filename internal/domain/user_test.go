package domain

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: User{
				Firstname: "John",
				Lastname:  "Doe",
				Login:     "johndoe",
				Password:  "password123",
			},
			wantErr: false,
		},
		{
			name: "empty firstname",
			user: User{
				Firstname: "",
				Lastname:  "Doe",
				Login:     "johndoe",
				Password:  "password123",
			},
			wantErr: true,
			errMsg:  "firstname is required",
		},
		{
			name: "whitespace firstname",
			user: User{
				Firstname: "   ",
				Lastname:  "Doe",
				Login:     "johndoe",
				Password:  "password123",
			},
			wantErr: true,
			errMsg:  "firstname is required",
		},
		{
			name: "empty lastname",
			user: User{
				Firstname: "John",
				Lastname:  "",
				Login:     "johndoe",
				Password:  "password123",
			},
			wantErr: true,
			errMsg:  "lastname is required",
		},
		{
			name: "empty login",
			user: User{
				Firstname: "John",
				Lastname:  "Doe",
				Login:     "",
				Password:  "password123",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
		{
			name: "login too short",
			user: User{
				Firstname: "John",
				Lastname:  "Doe",
				Login:     "ab",
				Password:  "password123",
			},
			wantErr: true,
			errMsg:  "login must be at least 3 characters",
		},
		{
			name: "empty password",
			user: User{
				Firstname: "John",
				Lastname:  "Doe",
				Login:     "johndoe",
				Password:  "",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "password too short",
			user: User{
				Firstname: "John",
				Lastname:  "Doe",
				Login:     "johndoe",
				Password:  "pass",
			},
			wantErr: true,
			errMsg:  "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}
