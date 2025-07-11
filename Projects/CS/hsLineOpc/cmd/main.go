package main

import (
	"context"
	"fmt"
	"hsLineOpc/api"
	"hsLineOpc/internal/handler"
	"hsLineOpc/pkg/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger.SetupLogging(slog.LevelDebug)

	err := godotenv.Load()
	if err != nil {
		slog.Error(fmt.Sprintf(".env file load error: %s", err.Error()))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startHs(ctx)
}

func startHs(ctx context.Context) {
	tsServ := api.NewTsClient(os.Getenv("TS_SERVER_CONN"))
	err := tsServ.Client.Connect(context.Background())
	if err != nil {
		log.Fatalf("TS server connect error: %s", err.Error())
	}
	defer tsServ.Client.Close()
	slog.Info("TS client connected")

	tsServ.SubscribeTs()
	slog.Info("TS server tags subscribed")

	opcClient := api.NewClient(os.Getenv("OPC_SERVER_CONN"))
	err = opcClient.Connect(context.Background())
	if err != nil {
		log.Fatalf("CS server connect error: %s", err.Error())
	}
	defer opcClient.Close()
	slog.Info("HS client connected")

	controlSys := handler.NewControlSystem(opcClient)
	go func() {
		ctxControl, cancel := context.WithCancel(context.Background())
		for {
			if tsServ.Start && !controlSys.IsActive && controlSys.IsDefault {
				slog.Info("Starting HS line")
				controlSys.Start(ctxControl)
			}

			if tsServ.Stop && controlSys.IsActive {
				err := controlSys.Stop()
				if err != nil {
					slog.Error(err.Error())
					continue
				}

				slog.Info("Stopping HS line")
				cancel()
				ctxControl, cancel = context.WithCancel(context.Background())
			}

			if tsServ.BackToStart && !controlSys.IsActive && !controlSys.IsDefault {
				slog.Info("Moving to default HS line")
				controlSys.Default()
			}

			time.Sleep(time.Millisecond * 10)
		}
	}()
	slog.Info("listening for TS server tags")

	<-ctx.Done()
	slog.Info("Shutting down control system gracefully")

	controlSys.Stop()
}
