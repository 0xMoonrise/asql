package server

import (
	"github.com/0xMoonrise/asql/internal/kernel"
	"github.com/gin-gonic/gin"
)

type app struct {
	kernel *kernel.Kernel
}

func NewServer() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies([]string{"172.18.0.0/16"})
	core := app{
		kernel: kernel.NewKernel("db.json"),
	}

	router.GET("/", core.root)
	router.GET("/check", core.health)
	router.POST("/run", core.run)

	return router
}
