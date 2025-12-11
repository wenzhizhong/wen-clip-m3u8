package routers

import (
	"clipM3u8Media/backend/apps/controllers"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	localFile := &controllers.LocalFileController{}
	router.GET("/play-local-file", localFile.PlayLocalFile)

}
