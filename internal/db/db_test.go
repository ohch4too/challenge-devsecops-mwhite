package db

import (
	"challenge/internal/domain"
	"os"
	"testing"

	"gorm.io/gorm"
)

func TestSqliteConnector(t *testing.T) {
	originalConn := Conn
	defer func() {
		Conn = originalConn
		os.Remove("./test.db")
	}()

	if err := sqliteConnector(); err != nil {
		t.Errorf("sqliteConnector() error = %v", err)
	}

	if Conn == nil {
		t.Error("sqliteConnector() did not set Conn")
	}
}

func TestPostgresConnection(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		user     string
		password string
		dbname   string
		wantErr  bool
	}{
		{
			name:     "missing password",
			host:     "localhost",
			user:     "testuser",
			password: "",
			dbname:   "testdb",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := postgresConnection(tt.host, tt.user, tt.password, tt.dbname); (err != nil) != tt.wantErr {
				t.Errorf("postgresConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		user     string
		password string
		dbname   string
		wantErr  bool
	}{
		{
			name:     "sqlite initialization",
			host:     "",
			user:     "",
			password: "",
			dbname:   "",
			wantErr:  false,
		},
		{
			name:     "postgres with missing password",
			host:     "localhost",
			user:     "testuser",
			password: "",
			dbname:   "testdb",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConn := Conn
			defer func() {
				Conn = originalConn
				if tt.host == "" {
					os.Remove("./test.db")
				}
			}()

			if err := Initialize(tt.host, tt.user, tt.password, tt.dbname); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && Conn == nil {
				t.Error("Initialize() did not set Conn")
			}
		})
	}
}

func TestSeedAdminUser(t *testing.T) {
	tests := []struct {
		name          string
		adminPassword string
		setupDB       func(*gorm.DB)
		wantErr       bool
		checkAdmin    bool
	}{
		{
			name:          "seed admin with custom password",
			adminPassword: "custompass",
			setupDB:       func(db *gorm.DB) {},
			wantErr:       false,
			checkAdmin:    true,
		},
		{
			name:          "seed admin with default password",
			adminPassword: "",
			setupDB:       func(db *gorm.DB) {},
			wantErr:       false,
			checkAdmin:    true,
		},
		{
			name:          "skip seeding when users exist",
			adminPassword: "testpass",
			setupDB: func(db *gorm.DB) {
				user := domain.User{
					Firstname: "Existing",
					Lastname:  "User",
					Login:     "existing",
					Password:  "password",
				}
				db.Create(&user)
			},
			wantErr:    false,
			checkAdmin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConn := Conn
			defer func() {
				Conn = originalConn
				os.Remove("./test.db")
			}()

			if err := Initialize("", "", "", ""); err != nil {
				t.Fatalf("Initialize() error = %v", err)
			}

			if tt.setupDB != nil {
				tt.setupDB(Conn)
			}

			if err := SeedAdminUser(tt.adminPassword); (err != nil) != tt.wantErr {
				t.Errorf("SeedAdminUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkAdmin {
				var admin domain.User
				if err := Conn.Where("login = ?", "admin").First(&admin).Error; err != nil {
					t.Error("SeedAdminUser() did not create admin user")
				}
				if admin.Firstname != "Admin" {
					t.Errorf("SeedAdminUser() admin firstname = %v, want Admin", admin.Firstname)
				}
			}
		})
	}
}
