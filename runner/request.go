package runner

import "github.com/yunhanshu-net/sdk-go/model/request"

type Request struct {
	Request         *request.Request `json:"request"`
	TransportConfig *TransportConfig `json:"transport_config"`
}
