package main

import (
	"hsLineOpc/api"
	"hsLineOpc/internal/handler"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

///go:generate go tool oapi-codegen -config ../api/config.yaml ../api/swagger.yaml

func main() {
	godotenv.Load()

	cfg := handler.Config{
		Port: os.Getenv("SERVER_PORT"),
	}

	server := handler.NewServer()
	mux := http.NewServeMux()
	h := api.HandlerFromMux(server, mux)

	srv := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Port,
		Handler: h,
	}

	slog.Info("Starting server on address " + srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen and serve: %v", err)
	}
}
