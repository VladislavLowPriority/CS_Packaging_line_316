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
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger.SetupLogging(slog.LevelDebug)

	err := godotenv.Load()
	if err != nil {
		slog.Error(fmt.Sprintf(".env file load error: %s", err.Error()))
	}

	tsServ := api.NewTsClient()

	err = tsServ.Client.Connect(context.Background())
	if err != nil {
		log.Fatalf("TS server connect error: %s", err.Error())
	}
	defer tsServ.Client.Close()

	tsServ.SubscribeTs()
	slog.Info("TS server tags subscribed")

	opcConn := os.Getenv("OPC_SERVER_IP") + ":" + os.Getenv("OPC_SERVER_PORT")
	opcClient := api.NewClient(opcConn)

	err = opcClient.Connect(context.Background())
	if err != nil {
		log.Fatalf("CS server connect error: %s", err.Error())
	}
	defer opcClient.Close()

	controlSys := handler.NewControlSystem(opcClient)
	ctx, cancel := context.WithCancel(context.Background())
	for {
		if tsServ.Start && !controlSys.IsActive && controlSys.IsDefault {
			controlSys.Start(ctx)
		}

		if tsServ.Stop && controlSys.IsActive {
			err := controlSys.Stop()
			if err != nil {
				slog.Error(err.Error())
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
