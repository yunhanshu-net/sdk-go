package request

import "encoding/json"

type OnApiCreated struct {
}

type AfterApiDelete struct {
}

type BeforeApiDelete struct {
}

type AfterApiDeleted struct {
}

type BeforeRunnerClose struct {
}

type AfterRunnerClose struct {
}

type OnVersionChange struct {
}

type OnInputFuzzy struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OnInputValidate struct {
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
