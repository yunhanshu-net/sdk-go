package response

type JSON interface {
	Builder
}

func (r *Data) JSON(data interface{}) error {
	r.DataType = DataTypeJSON
	r.Body = data
	r.StatusCode = successCode
	r.Msg = successMsg
	return nil
}
