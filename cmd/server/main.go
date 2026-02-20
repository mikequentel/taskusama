package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mikequentel/taskusama/internal/httpapi"
	"github.com/mikequentel/taskusama/internal/store/memory"
)

func main() {
	issues := memory.New()
	api := httpapi.New(issues)

	r := gin.New()
	r.Use(gin.Recovery())

	// templates + static
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "web/static")

	api.RegisterRoutes(r)

	_ = r.Run(":8080")
}
