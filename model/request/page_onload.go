package request

import (
	"encoding/json"
	"fmt"
)

type OnApiCreated struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeApiDelete struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type AfterApiDeleted struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeRunnerClose struct {
}

type AfterRunnerClose struct {
}

type Change struct {
	Method string `json:"method"`
	Router string `json:"router"`
	Type   string `json:"type"`
}

func (c *Change) String() string {
	return fmt.Sprintf(`{"method": "%s", "router": "%s","type","%s"}`, c.Method, c.Router, c.Type)
}

type OnVersionChange struct {
	Change []Change `json:"change"`
}

type OnInputFuzzy struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OnInputValidate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OnTableDeleteRows struct {
	Ids []string `json:"ids"`
}

type OnTableUpdateRow struct {
	Ids []string `json:"ids"`
}

type OnTableSearch struct {
	Cond map[string]string `json:"cond"`
}

type Callback struct {
	Method string      `json:"method"`
	Router string      `json:"router"`
	Type   string      `json:"type"`
	Body   interface{} `json:"body"`
}

func (c *Callback) DecodeData(el interface{}) error {
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

//func (c *Callback) BindData(req interface{}) error {
//	err := json.Unmarshal([]byte(c.Body), &req)
//	if err != nil {
//		return err
//	}
//	return nil
//}
