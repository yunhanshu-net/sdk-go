package response

type Files interface {
	AddFile(localFile *LocalFile) Files
	Build() error
}

type fileData struct {
	Title string       `json:"title"`
	Files []*LocalFile `json:"files"`
}

func (f *fileData) DataType() DataType {
	return DataTypeFiles
}
