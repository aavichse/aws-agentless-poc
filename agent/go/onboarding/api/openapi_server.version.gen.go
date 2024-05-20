// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package api

import (
	"github.com/gin-gonic/gin"
)

// PostVersionHandshake operation middleware
func (siw *ServerInterfaceWrapper) PostVersionHandshake(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostVersionHandshake(c)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlersVersion(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptionsVersion(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptionsVersion(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/version/version-handshake", wrapper.PostVersionHandshake)
}
