package request

import "encoding/json"

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
}

type OnTableUpdateRow struct {
}

type OnTableSearch struct {
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
