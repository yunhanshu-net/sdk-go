package api

import (
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
)

func NewResponseParams(el interface{}, renderType string) (interface{}, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	switch renderType {
	case response.RenderTypeForm:
		return NewFormResponseParams(el)
	case response.RenderTypeTable:
		return NewTableResponseParams(el)
	default:
		return NewFormResponseParams(el)
	}
}
