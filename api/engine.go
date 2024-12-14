package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(debug bool) *gin.Engine {
	if debug {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.Default()
	engine.Use(gin.Recovery())

	return mountRoutes(engine)
}

func mountRoutes(e *gin.Engine) *gin.Engine {
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "YueWen.store.pusher",
			"status":   "OK",
		})
	})
	
	return e
}