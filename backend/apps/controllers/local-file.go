package controllers

import (
	"clipM3u8Media/backend/apps/common/controller"
	"clipM3u8Media/backend/apps/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type LocalFileController struct {
	controller.BaseController
	LocalFileService services.LocalFileService
}

// 播放本地文件
func (c *LocalFileController) PlayLocalFile(ctx *gin.Context) {
	path := ctx.Query("path")
	content, err := c.LocalFileService.PlayLocalFile(path)

	code := 200
	msg := "success"
	if err != nil {
		code = 500
		msg = err.Error()

		fmt.Println("PlayLocalFile(): error：", msg)
	}
	// 跨域处理

	// 设置响应头
	ctx.Header("Content-Type", "video/mp4")
	ctx.Header("Accept-Ranges", "bytes") // 支持断点续传

	if code == 200 {
		//输出二进制
		ctx.Data(code, "video/mp4", *content)
	} else {
		ctx.Data(code, "video/mp4", nil)
	}
}
