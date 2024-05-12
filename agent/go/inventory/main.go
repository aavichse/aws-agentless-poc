package main

import (
	logger "agentless/infra/log"
	"flag"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultLogPath       = "/var/log/agentless/inventory.log"
	defaultFetchInterval = 10
	defaultMaxConcurrency = 3
)

func main() {
	fetchInterval := flag.Int("interval", defaultFetchInterval, "Interval of fetching full update (in seconds)")
	logPath := flag.String("logPath", defaultLogPath, "Path to the log file")
	logDebug := flag.Bool("debug", false, "Set log level to debug")
	flag.Parse()

	logger.InitLogger(*logPath, *logDebug)
	logger.Log.Info("Starting")

	router := gin.Default()
	inventoryService := NewInventoryService()

	RegisterHandlers(router, inventoryService)

	go func() {
		if err := router.Run("0.0.0.0:8080"); err != nil {
			panic(err)
		}
	}()

	regions := []string{
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
	}
	interval := time.Duration(*fetchInterval) * time.Second
	fetcher := NewAWSFetcher(interval, defaultMaxConcurrency, regions, inventoryService)

	fetcher.Start()

	<-fetcher.done
}
