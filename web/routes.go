package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// OpenCORSMiddleware controls the cross origin policies.
func OpenCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Handler is an interface to the HTTP handler functions.
type Handler interface {
	CreateRecord(c *gin.Context)
	SearchRecords(c *gin.Context)
}

// SetRoutes sets all of the appropriate routes to handlers for the application
func SetRoutes(engine *gin.Engine, h Handler) error {
	api := engine.Group("/")

	api.Use(OpenCORSMiddleware())

	api.POST("/:project", h.CreateRecord)
	api.GET("/:project", h.SearchRecords)

	return nil
}
