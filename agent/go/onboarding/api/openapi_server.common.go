package api

import (
	model "agentless/infra/model/operations"

	"github.com/gin-gonic/gin"
)

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Endpoint to retrieve internal configuration options
	// (GET /v1/operations/internal/config_metadata)
	GetV1OperationsInternalConfigMetadata(c *gin.Context)
	// Endpoint to set internal services configuration
	// (POST /v1/operations/internal/internal_config)
	PostV1OperationsInternalInternalConfig(c *gin.Context)
	// Endpoint to retrieve components status
	// (GET /v1/operations/health)
	GetV1OperationsHealth(c *gin.Context)
	// Endpoint to retrieve components metrics
	// (GET /v1/operations/metrics)
	GetV1OperationsMetrics(c *gin.Context)
	// Endpoint to retrieve integration environment details. Should fit in one response page.
	// (GET /v1/operations/env-info)
	GetV1OperationsEnvInfo(c *gin.Context)
	// Endpoint to retrieve integration environment logical units details. This endpoint support pagination for big amount of units
	// (GET /v1/operations/env-units-list)
	GetV1OperationsEnvUnitsList(c *gin.Context, params model.GetV1OperationsEnvUnitsListParams)
	// Endpoint to retrieve component logs
	// (GET /v1/operations/log/download)
	GetV1OperationsLogDownload(c *gin.Context, params model.GetV1OperationsLogDownloadParams)
	// Endpoint for instructing component to gather logs
	// (POST /v1/operations/log/start)
	PostV1OperationsLogStart(c *gin.Context)
	// Endpoint to query component on logs gathering status
	// (GET /v1/operations/log/status)
	GetV1OperationsLogStatus(c *gin.Context, params model.GetV1OperationsLogStatusParams)
	// Endpoint to abort logs gathering
	// (GET /v1/operations/log/stop)
	GetV1OperationsLogStop(c *gin.Context, params model.GetV1OperationsLogStopParams)
	// Endpoint to onboard connector
	// (POST /v1/operations/onboarding)
	PostV1OperationsOnboarding(c *gin.Context)
	// Endpoint to get onboarding status
	// (POST /v1/operations/status)
	PostV1OperationsStatus(c *gin.Context)
	// Perform version handshake
	// (POST /version-handshake)
	PostVersionHandshake(c *gin.Context)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	RegisterHandlersWithOptionsConfig(router, si, options)
	RegisterHandlersWithOptionsHealth(router, si, options)
	RegisterHandlersWithOptionsInfo(router, si, options)
	RegisterHandlersWithOptionsLog(router, si, options)
	RegisterHandlersWithOptionsOnboard(router, si, options)
	RegisterHandlersWithOptionsVersion(router, si, options)
}
