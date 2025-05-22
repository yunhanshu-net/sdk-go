package api

import (
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
)

func NewRequestParams(el interface{}, renderType string) (interface{}, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	switch renderType {
	case response.RenderTypeTable:
		return NewTableRequestParams(el)
	case response.RenderTypeForm:
		return NewFormRequestParams(el, renderType)
	default:
		return NewFormRequestParams(el, renderType)
	}
}
