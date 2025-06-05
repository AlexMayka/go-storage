package main

import (
	"fmt"
	_ "go-storage/cmd/api/docs"
	"go-storage/internal/config"
	"go-storage/internal/delivery/http"
	"go-storage/pkg/db"
	"go-storage/pkg/logger"
	"log"
)

// @title       Go-Storage
// @version     1.0
// @description This project is being developed as a cloud storage
// @host        localhost:8080
// @BasePath    /api/v1
func main() {
	cfg := config.NewConfig()

	saveToFile := cfg.App.LogToFile == "true"
	logging, errLgn := logger.InitLog(saveToFile, cfg.App.LogLevel)
	if errLgn != nil {
		log.Fatalf("InitLog init failed: %v", errLgn)
		return
	}

	database, errDb := db.InitDB(
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Name,
	)
	if errDb != nil {
		logging.Error("InitDB init failed", "error", errDb)
		return
	}

	r := http.Router(logging, database)
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)

	logging.Info("Run server", "addr", addr)
	errRun := r.Run(addr)

	if errRun != nil {
		logging.Error("Start server failed", "error", errRun)
	}
}
