package response

import "net/http"

type JSON interface {
	Builder
}

func (r *Data) JSON(data interface{}) error {
	r.RenderType = RenderTypeJSON
	r.Body = data
	r.StatusCode = http.StatusOK
	return nil
}
