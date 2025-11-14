package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	originalHost := os.Getenv("POSTGRES_HOST")
	originalUser := os.Getenv("POSTGRES_USER")
	originalDB := os.Getenv("POSTGRES_DB")
	originalPassword := os.Getenv("POSTGRES_PASSWORD")
	originalAdminPassword := os.Getenv("ADMIN_PASSWORD")
	originalTLSCert := os.Getenv("TLS_CERT_FILE")
	originalTLSKey := os.Getenv("TLS_KEY_FILE")

	defer func() {
		os.Setenv("POSTGRES_HOST", originalHost)
		os.Setenv("POSTGRES_USER", originalUser)
		os.Setenv("POSTGRES_DB", originalDB)
		os.Setenv("POSTGRES_PASSWORD", originalPassword)
		os.Setenv("ADMIN_PASSWORD", originalAdminPassword)
		os.Setenv("TLS_CERT_FILE", originalTLSCert)
		os.Setenv("TLS_KEY_FILE", originalTLSKey)
	}()

	tests := []struct {
		name       string
		setupEnv   func()
		wantErr    bool
		wantDBHost string
		wantDBUser string
		wantDBName string
	}{
		{
			name: "default values with no postgres host",
			setupEnv: func() {
				os.Unsetenv("POSTGRES_HOST")
				os.Unsetenv("POSTGRES_USER")
				os.Unsetenv("POSTGRES_DB")
				os.Unsetenv("POSTGRES_PASSWORD")
			},
			wantErr:    false,
			wantDBHost: "",
			wantDBUser: "challenge",
			wantDBName: "challenge",
		},
		{
			name: "custom postgres values",
			setupEnv: func() {
				os.Setenv("POSTGRES_HOST", "testhost")
				os.Setenv("POSTGRES_USER", "testuser")
				os.Setenv("POSTGRES_DB", "testdb")
				os.Setenv("POSTGRES_PASSWORD", "testpass")
			},
			wantErr:    false,
			wantDBHost: "testhost",
			wantDBUser: "testuser",
			wantDBName: "testdb",
		},
		{
			name: "postgres host set but no password",
			setupEnv: func() {
				os.Setenv("POSTGRES_HOST", "localhost")
				os.Unsetenv("POSTGRES_PASSWORD")
			},
			wantErr: true,
		},
		{
			name: "all environment variables set",
			setupEnv: func() {
				os.Setenv("POSTGRES_HOST", "localhost")
				os.Setenv("POSTGRES_USER", "user")
				os.Setenv("POSTGRES_DB", "db")
				os.Setenv("POSTGRES_PASSWORD", "pass")
				os.Setenv("ADMIN_PASSWORD", "adminpass")
				os.Setenv("TLS_CERT_FILE", "/path/to/cert")
				os.Setenv("TLS_KEY_FILE", "/path/to/key")
			},
			wantErr:    false,
			wantDBHost: "localhost",
			wantDBUser: "user",
			wantDBName: "db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg.DBHost != tt.wantDBHost {
					t.Errorf("Load() DBHost = %v, want %v", cfg.DBHost, tt.wantDBHost)
				}
				if cfg.DBUser != tt.wantDBUser {
					t.Errorf("Load() DBUser = %v, want %v", cfg.DBUser, tt.wantDBUser)
				}
				if cfg.DBName != tt.wantDBName {
					t.Errorf("Load() DBName = %v, want %v", cfg.DBName, tt.wantDBName)
				}
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		want         string
	}{
		{
			name:         "env var not set, use default",
			key:          "TEST_VAR_1",
			defaultValue: "default",
			want:         "default",
			setEnv:       false,
		},
		{
			name:         "env var set, use env value",
			key:          "TEST_VAR_2",
			defaultValue: "default",
			envValue:     "custom",
			setEnv:       true,
			want:         "custom",
		},
		{
			name:         "env var set to empty, use default",
			key:          "TEST_VAR_3",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnvOrDefault(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
