package response

type fileData struct {
	Title string  `json:"title"`
	Files []*File `json:"files"`
}

func (f *fileData) DataType() DataType {
	return DataTypeFiles
}
