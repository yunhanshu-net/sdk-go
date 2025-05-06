package callback

import "encoding/json"

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
