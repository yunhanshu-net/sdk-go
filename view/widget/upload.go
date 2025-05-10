package widget

// UploadWidget 文件上传组件
type UploadWidget struct {
	// 组件类型，可选值：upload(普通上传)、image-upload(图片上传)、drag-upload(拖拽上传)
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
}

func (w *UploadWidget) GetValueType() string {
	return w.Type
}

func (w *UploadWidget) GetWidgetType() string {
	return w.Widget
}
