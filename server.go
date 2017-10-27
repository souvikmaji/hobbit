package main

import (
	"github.com/gin-gonic/gin"
)

func StartServer(cfg *Config) {
	r := gin.Default()
	setupRoutes(r)
	r.Run(cfg.Addr())
}
