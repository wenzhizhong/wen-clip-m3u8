package goApi

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Runtime struct {
	Ctx *context.Context
}

func (a *Runtime) MessageDialog(msg, title string) (string, error) {
	if title == "" {
		title = "提示"
	}
	return runtime.MessageDialog(*a.Ctx, runtime.MessageDialogOptions{
		Title:         title,
		Message:       msg,
		Buttons:       []string{"取消", "确定"},
		DefaultButton: "确定",
	})
}

// OpenFileDialog
func (a *Runtime) OpenFileDialog(optionJson map[string]interface{}) (string, error) {

	fmt.Println(optionJson)
	option := runtime.OpenDialogOptions{
		Title: "选择文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "所有文件",
				Pattern:     "*.*",
			},
		},
	}
	if optionJson != nil {
		if optionJson["Title"] != nil {
			option.Title = optionJson["Title"].(string)
		}
		if optionJson["Filters"] != nil {
			option.Filters = make([]runtime.FileFilter, 0)
			for _, filter := range optionJson["Filters"].([]interface{}) {
				option.Filters = append(option.Filters, runtime.FileFilter{
					DisplayName: filter.(map[string]interface{})["DisplayName"].(string),
					Pattern:     filter.(map[string]interface{})["Pattern"].(string),
				})
			}
		}
	}

	return runtime.OpenFileDialog(*a.Ctx, option)
}
