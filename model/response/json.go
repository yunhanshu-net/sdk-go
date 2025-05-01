package response

import "net/http"

type JSON interface {
	Builder
}

func (r *Data) JSON(data interface{}) error {
	r.DataType = DataTypeJSON
	r.Body = data
	r.StatusCode = http.StatusOK
	return nil
}
