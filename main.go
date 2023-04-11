package main

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherm5/memcachedGorm/routes"
)

func main() {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(nil))

	routes.Users(router)

	router.Run(":8085")
}
