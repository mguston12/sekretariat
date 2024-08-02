package boot

import (
	"log"
	"net/http"

	"sekretariat/internal/config"

	"github.com/jmoiron/sqlx"

	sekretariatData "sekretariat/internal/data/sekretariat"
	sekretariatServer "sekretariat/internal/delivery/http"
	sekretariatHandler "sekretariat/internal/delivery/http/sekretariat"
	sekretariatService "sekretariat/internal/service/sekretariat"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
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

	// Diganti dengan domain yang anda buat
	sd := sekretariatData.New(db)
	ss := sekretariatService.New(sd)
	sh := sekretariatHandler.New(ss)

	server := sekretariatServer.Server{
		Sekretariat: sh,
	}

	if err := server.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
