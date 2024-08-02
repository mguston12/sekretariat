package boot

import (
	"log"
	"net/http"
	"os"

	"sekretariat/internal/config"

	"github.com/jmoiron/sqlx"

	sekretariatData "sekretariat/internal/data/sekretariat"
	sekretariatServer "sekretariat/internal/delivery/http"
	sekretariatHandler "sekretariat/internal/delivery/http/sekretariat"
	sekretariatService "sekretariat/internal/service/sekretariat"
)

// HTTP will load configuration, do dependency injection, and then start the HTTP server
func HTTP() error {
	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg := config.Get()
	// Open MySQL DB Connection
	db, err := sqlx.Open("mysql", cfg.Database.Master)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	// Dependency Injection
	sd := sekretariatData.New(db)
	ss := sekretariatService.New(sd)
	sh := sekretariatHandler.New(ss)

	server := sekretariatServer.Server{
		Sekretariat: sh,
	}

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	if err := server.Serve(":" + port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
