package main

import (
	"challenge/internal/api"
	"challenge/internal/config"
	"challenge/internal/db"
	"fmt"
	"os"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	if err := db.Initialize(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Seed admin user
	if err := db.SeedAdminUser(cfg.AdminPassword); err != nil {
		fmt.Printf("Failed to seed admin user: %v\n", err)
		os.Exit(3)
	}

	// Start server
	router := api.SetupRouter(db.Conn)

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		router.RunTLS(":10000", cfg.TLSCertFile, cfg.TLSKeyFile)
	} else {
		router.Run(":10000")
	}
}
