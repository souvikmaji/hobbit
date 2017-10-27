package main

import "github.com/gin-gonic/gin"

func setupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
    c.Request.Header.Add("Content-Type", "text/html")
		c.String(200, LAYOUT)
	})
}
