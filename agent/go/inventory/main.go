package main

import (
	logger "agentless/infra/log"
	"flag"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultLogPath        = "/var/log/guardicore/inventory.log"
	defaultFetchInterval  = 60
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
		if err := router.Run("0.0.0.0:8082"); err != nil {
			panic(err)
		}
	}()

	regions := []string{
		"us-east-1", // US East (N. Virginia)
		"us-east-2", // US East (Ohio)
		"us-west-1", // US West (N. California)
		"us-west-2", // US West (Oregon)
		// "ca-central-1",   // Canada (Central)
		// "ca-west-1",      // Canada (West)
		"eu-north-1",   // EU (Stockholm)
		"eu-west-3",    // EU (Paris)
		"eu-west-2",    // EU (London)
		"eu-west-1",    // EU (Ireland)
		"eu-central-1", // EU (Frankfurt)
		// "eu-south-1",   // EU (Milan)
		// "ap-south-1",     // Asia Pacific (Mumbai)
		"ap-northeast-1", // Asia Pacific (Tokyo)
		"ap-northeast-2", // Asia Pacific (Seoul)
		"ap-northeast-3", // Asia Pacific (Osaka-Local)
		"ap-southeast-1", // Asia Pacific (Singapore)
		"ap-southeast-2", // Asia Pacific (Sydney)
		// "ap-southeast-3", // Asia Pacific (Jakarta)
		// "ap-east-1", // Asia Pacific (Hong Kong) SAR
		// "sa-east-1", // South America (São Paulo)
		// "cn-north-1",     // China (Beijing)
		// "cn-northwest-1", // China (Ningxia)
		// "us-gov-east-1",  // GovCloud (US-East)
		// "us-gov-west-1",  // GovCloud (US-West)
		// "us-isob-east-1", // AWS Secret Region (US ISOB East Ohio)
		// "us-iso-east-1",  // AWS Top Secret-East Region (US ISO East Virginia)
		// "us-iso-west-1",  // AWS Top Secret-West Region (US ISO West Colorado)
		// "me-south-1",     // Middle East (Bahrain)
		// "af-south-1",     // Africa (Cape Town)
		// "me-central-1",   // Middle East (United Arab Emirates)
		// "eu-south-2",     // EU (Spain)
		// "eu-central-2",   // EU (Zurich)
		// "ap-south-2",     // Asia Pacific (Hyderabad)
		// "ap-southeast-4", // Asia Pacific (Melbourne)
		// "il-central-1",   // Israel (Tel Aviv)	}
	}

	interval := time.Duration(*fetchInterval) * time.Second
	fetcher := NewAWSFetcher(interval, defaultMaxConcurrency, regions, inventoryService)

	fetcher.Start()

	<-fetcher.done
}
