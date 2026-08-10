package goApi

import (
	"clipM3u8Media/goApi/common"
	"context"
)

type GetConstant struct {
	Ctx *context.Context
}

func (g *GetConstant) GetAppMainWindowName() string {
	return common.AppMainWindowName
}

func (g *GetConstant) GetAppPreUploadWindowName() string {
	return common.AppPreUploadWindowName
}
