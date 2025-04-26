package user

import "github.com/gin-gonic/gin"

func MeHandler(c *gin.Context) {
	c.String(200, "Là mình đây1 ???")
}