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

	router.LoadHTMLGlob("templates/*")
	router.Static("/js", "/opt/lib/htmx")

	router.GET("/", core.root)
	router.GET("/check", core.health)
	router.POST("/lexer", core.lexer)

	return router
}
