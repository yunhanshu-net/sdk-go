package widget

import "github.com/yunhanshu-net/sdk-go/pkg/tagx"

// FileWidget 文件上传组件
type FileWidget struct {
	// 组件类型，固定为file
	Widget string `json:"widget"`
	// 数据类型，一般为file或string(文件路径)
	Type string `json:"type"`
	// 接受的文件类型，如：.jpg,.png,.pdf
	Accept string `json:"accept,omitempty"`
	// 是否支持多文件上传
	Multiple bool `json:"multiple,omitempty"`
	// 文件大小限制，单位KB
	MaxSize int `json:"max_size,omitempty"`
	// 上传文件数量限制
	Limit int `json:"limit,omitempty"`
	// 占位符/提示文本
	Placeholder string `json:"placeholder,omitempty"`
	// 是否自动上传
	AutoUpload bool `json:"auto_upload,omitempty"`
	// 上传接口地址
	Action string `json:"action,omitempty"`
	// 列表类型，可选值：text(文本)、picture(图片)、picture-card(卡片)
	ListType string `json:"list_type,omitempty"`
	// 是否拖拽上传
	Drag bool `json:"drag,omitempty"`
	// 上传按钮文字
	ButtonText string `json:"button_text,omitempty"`
	// 提示文字
	Tip string `json:"tip,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
}

// newFileWidget 创建文件上传组件
func newFileWidget(info *tagx.FieldInfo) (Widget, error) {
	file := &FileWidget{
		Widget: WidgetFile,
		Type:   TypeFile,
	}

	tag := info.Tags
	if tag["accept"] != "" {
		file.Accept = tag["accept"]
	}

	if tag["multiple"] != "" {
		if tag["multiple"] == "true" {
			file.Multiple = true
		}
	}

	if tag["max_size"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		file.MaxSize = 0
	}

	if tag["limit"] != "" {
		// 这里可以添加字符串转整数的逻辑，简化处理
		file.Limit = 0
	}

	if tag["placeholder"] != "" {
		file.Placeholder = tag["placeholder"]
	}

	if tag["auto_upload"] != "" {
		if tag["auto_upload"] == "true" {
			file.AutoUpload = true
		}
	}

	if tag["action"] != "" {
		file.Action = tag["action"]
	}

	if tag["list_type"] != "" {
		file.ListType = tag["list_type"]
	}

	if tag["drag"] != "" {
		if tag["drag"] == "true" {
			file.Drag = true
		}
	}

	if tag["button_text"] != "" {
		file.ButtonText = tag["button_text"]
	}

	if tag["tip"] != "" {
		file.Tip = tag["tip"]
	}

	return file, nil
}

func (w *FileWidget) GetValueType() string {
	return w.Type
}

func (w *FileWidget) GetWidgetType() string {
	return w.Widget
}
