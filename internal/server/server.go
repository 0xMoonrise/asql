package server

import (
	"github.com/gin-gonic/gin"
)

type app struct {
	// dependecies
}

func NewServer() *gin.Engine {
	router := gin.Default()

	core := app{}

	router.GET("/", core.root)
	router.GET("/check", core.health)
	router.POST("/run", core.run)

	return router
}
