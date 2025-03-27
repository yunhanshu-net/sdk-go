package response

type LocalFile struct {
	Title string `json:"title"` //2020年语文真题
	Path  string `json:"path"`  //文件地址
}

type files struct {
	response   *Response
	Title      string       //	2020年-2025年考研历年真题汇总
	localFiles []*LocalFile //文件列表
	value      interface{}
}

func (r *Response) Files(localFiles []*LocalFile, title ...string) Files {
	f := files{
		response:   r,
		localFiles: localFiles,
	}
	if len(title) > 0 {
		f.Title = title[0]
	}
	return &f
}

func (f *files) AddFile(localFile *LocalFile) Files {
	f.localFiles = append(f.localFiles, localFile)
	return f
}
func (f *files) Build() error {
	f.response.DataType = DataTypeFiles
	f.response.data = append(f.response.data, &fileData{
		Title: f.Title,
		Files: f.localFiles,
	})
	return nil
}
