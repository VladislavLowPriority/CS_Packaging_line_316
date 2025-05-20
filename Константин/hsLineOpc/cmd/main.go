package main

import (
	"context"
	"hsLineOpc/api"
	"hsLineOpc/internal/handler"
	"hsLineOpc/pkg/logger"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger.SetupLogging(slog.LevelDebug)

	tsServ := api.NewTsClient()
	tsServ.SubscribeTs()
	slog.Info("TS server tags subscribed")

	opcConn := os.Getenv("OPC_SERVER_IP") + ":" + os.Getenv("OPC_SERVER_PORT")
	opcClient := api.NewClient(opcConn)
	slog.Info("HS client created")

	controlSys := handler.NewControlSystem(opcClient)
	ctx, cancel := context.WithCancel(context.Background())
	for {
		if tsServ.Start && !controlSys.IsActive && controlSys.IsDefault {
			controlSys.Start(ctx)
		}

		if tsServ.Stop && controlSys.IsActive {
			err := controlSys.Stop()
			if err != nil {
				continue
			}

			cancel()
			ctx, cancel = context.WithCancel(context.Background())
		}

		if tsServ.BackToStart && !controlSys.IsActive && !controlSys.IsDefault {
			controlSys.Default()
		}

		time.Sleep(time.Millisecond * 10)
	}
}
