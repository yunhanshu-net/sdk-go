package runner

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/yunhanshu-net/pkg/constants"
	"github.com/yunhanshu-net/pkg/x/jsonx"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/usercall"
)

func (r *Runner) syscallCmd(cmd *cobra.Command, args []string) {

}

func (r *Runner) userCallCmd(cmd *cobra.Command, args []string) {
	req := &usercall.Request{}
	resp := &response.RunFunctionResp{}
	callType, err := cmd.Flags().GetString("type")
	if err != nil {
		writeJSON(resp)
		resp.Msg = err.Error()
		return
	}
	method, err := cmd.Flags().GetString("method")
	if err != nil {
		writeJSON(resp)
		resp.Msg = err.Error()
		return
	}
	router, err := cmd.Flags().GetString("router")
	if err != nil {
		writeJSON(resp)
		resp.Msg = err.Error()
		return
	}
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		writeJSON(resp)
		resp.Msg = err.Error()
		return
	}

	if file != "noBody" && file != "" {
		err := jsonx.UnmarshalFromFile(file, &req)
		if err != nil {
			resp.Msg = err.Error()
			writeJSON(resp)
			return
		}
	} else {
		req.Method = method
		req.Router = router
		req.Type = callType
	}
	traceId, err := cmd.Flags().GetString("trace_id")
	ctx := &Context{
		Context: context.WithValue(context.Background(), constants.TraceID, traceId),
	}

	err = r._callback(ctx, req, resp)
	if err != nil {
		resp.Msg = err.Error()
		writeJSON(resp)
		return
	}
	resp.Code = 0
	resp.Msg = "ok"
	writeJSON(resp)

}
func (r *Runner) apisCmd(cmd *cobra.Command, args []string) {
	apiList, err := r.getApiInfos()
	if err != nil {
		panic(err)
	}
	writeJSON(apiList)
}
