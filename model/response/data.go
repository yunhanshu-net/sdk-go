package response

type anyData interface {
	Build() error
	DataType() DataType
}
