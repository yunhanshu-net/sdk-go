package usercall

import (
	"encoding/json"
	"fmt"
)

type OnApiCreatedReq struct {
	//Method string `json:"method"`
	//Router string `json:"router"`
}

type OnApiUpdatedReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeApiDeleteReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type AfterApiDeletedReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeRunnerCloseReq struct {
}

type AfterRunnerCloseReq struct {
}

type Change struct {
	Method string `json:"method"`
	Router string `json:"router"`
	Type   string `json:"type"`
}

func (c *Change) String() string {
	return fmt.Sprintf(`{"method": "%s", "router": "%s","type","%s"}`, c.Method, c.Router, c.Type)
}

type OnVersionChangeReq struct {
	Change []Change `json:"change"`
}

type OnInputFuzzyReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OnInputValidateReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OnTableDeleteRowsReq struct {
	Ids []string `json:"ids"`
}

type OnTableUpdateRowReq struct {
	Ids []string `json:"ids"`
}

type OnTableSearchReq struct {
	Cond map[string]string `json:"cond"`
}
type Request struct {
	Method string      `json:"method"`
	Router string      `json:"router"`
	Type   string      `json:"type"`
	Body   interface{} `json:"body"`
}

func (c *Request) DecodeData(el interface{}) error {
	marshal, err := json.Marshal(c.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &el)
	if err != nil {
		return err
	}
	return nil
}

type Response struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}
