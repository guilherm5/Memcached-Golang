package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherm5/memcachedGorm/controllers"
)

func Users(c *gin.Engine) {
	c.POST("/users", controllers.UsersFunc())
	c.GET("/users", controllers.UsersFuncID())
}
