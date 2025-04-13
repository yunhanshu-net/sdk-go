package runner

import (
	"context"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"strings"
)

type routerInfo struct {
	Handel interface{}
	key    string
	Router string
	Method string
	Config *ApiConfig
}

func (r *routerInfo) IsDefaultRouter() bool {
	return strings.HasPrefix(strings.TrimPrefix(r.Router, "/"), "_")
}

func fmtKey(router string, method string) string {
	if !strings.HasPrefix(router, "/") {
		router = "/" + router
	}
	router = strings.TrimSuffix(router, "/")
	return router + "." + strings.ToUpper(method)
}

func (r *routerInfo) call(ctx context.Context, reqBody interface{}) (req *request.Request, resp *response.Data, err error) {

	meta, ok := handlerCacheMap[r.key]
	if !ok {
		h := &handlerMeta{}
		handlerCacheMap[r.key] = h
		meta = h
	}
	//metaVal, _ := handlerCache.LoadOrStore(r.key, &handlerMeta{})
	//meta := metaVal.(*handlerMeta)
	// 确保只初始化一次
	meta.once.Do(func() {
		meta.meta, meta.initError = buildRuntimeMeta(r.Handel)
	})
	if meta.initError != nil {
		return nil, nil, meta.initError
	}
	req = new(request.Request)
	resp = new(response.Data)
	ctx1 := &Context{Context: ctx}
	err = doCall(r.Method, meta.meta, ctx1, resp, reqBody)
	if err != nil {
		return nil, nil, err
	}
	return req, resp, nil

}
