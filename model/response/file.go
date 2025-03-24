package response

type File struct {
	Title string `json:"title"` //2020年语文真题
	Path  string `json:"path"`  //文件地址
}

type Files struct {
	response *Response
	Title    string  `json:"title"` //	2020年-2025年考研历年真题汇总
	Files    []*File `json:"files"` //文件列表
	value    interface{}
}

func (r *Response) Files(files []*File, title ...string) *Files {
	f := Files{
		response: r,
		Files:    files,
	}
	if len(title) > 0 {
		f.Title = title[0]
	}
	return &f
}

func (f *Files) AddFile(file *File) *Files {
	f.Files = append(f.Files, file)
	return f
}
func (f *Files) Build() error {
	f.response.DataType = DataTypeFiles
	f.response.data = append(f.response.data, &fileData{
		Title: f.Title,
		Files: f.Files,
	})
	return nil
}
