package main

import (
	logger "agentless/infra/log"
	"flag"

	"agentless/onboarding/api"

	"github.com/gin-gonic/gin"
)

const (
	defaultLogPath = "/var/log/agentless/inventory.log"

	// FIXME: Hard coded for the POC
	defaultConnectorID = "cf178df3-1d8c-46a5-86b7-974a941c4d80"
)

func main() {
	logPath := flag.String("logPath", defaultLogPath, "Path to the log file")
	logDebug := flag.Bool("debug", false, "Set log level to debug")
	flag.Parse()

	logger.InitLogger(*logPath, *logDebug)
	logger.Log.Info("Starting")

	router := gin.Default()
	onboardingService := NewOnboardingService(defaultConnectorID)

	api.RegisterHandlers(router, onboardingService)

	go func() {
		if err := router.Run("0.0.0.0:8080"); err != nil {
			panic(err)
		}
	}()

	select {}
}
