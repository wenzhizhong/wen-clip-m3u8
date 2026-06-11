package goApi

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Runtime struct {
	App *application.App
}

func (a *Runtime) MessageDialog(msg, title string) {
	if title == "" {
		title = "提示"
	}
	// return runtime.MessageDialog(*a.Ctx, runtime.MessageDialogOptions{
	// 	Title:         title,
	// 	Message:       msg,
	// 	Buttons:       []string{"取消", "确定"},
	// 	DefaultButton: "确定",
	// })
	a.App.Dialog.Info().
		SetTitle(title).
		SetMessage(msg).
		Show()
}

// OpenFileDialog
func (a *Runtime) OpenFileDialog(optionJson map[string]interface{}) (string, error) {

	fmt.Println(optionJson)
	// option := runtime.OpenDialogOptions{
	// 	Title: "选择文件",
	// 	Filters: []runtime.FileFilter{
	// 		{
	// 			DisplayName: "所有文件",
	// 			Pattern:     "*.*",
	// 		},
	// 	},
	// }
	// if optionJson != nil {
	// 	if optionJson["Title"] != nil {
	// 		option.Title = optionJson["Title"].(string)
	// 	}
	// 	if optionJson["Filters"] != nil {
	// 		option.Filters = make([]runtime.FileFilter, 0)
	// 		for _, filter := range optionJson["Filters"].([]interface{}) {
	// 			option.Filters = append(option.Filters, runtime.FileFilter{
	// 				DisplayName: filter.(map[string]interface{})["DisplayName"].(string),
	// 				Pattern:     filter.(map[string]interface{})["Pattern"].(string),
	// 			})
	// 		}
	// 	}
	// }

	// return runtime.OpenFileDialog(*a.Ctx, option)

	option := &application.OpenFileDialogOptions{
		Title: "选择文件",
		Filters: []application.FileFilter{
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
			option.Filters = make([]application.FileFilter, 0)
			for _, filter := range optionJson["Filters"].([]interface{}) {
				option.Filters = append(option.Filters, application.FileFilter{
					DisplayName: filter.(map[string]interface{})["DisplayName"].(string),
				})
			}
		}
	}
	return a.App.Dialog.OpenFileWithOptions(option).PromptForSingleSelection()
}
