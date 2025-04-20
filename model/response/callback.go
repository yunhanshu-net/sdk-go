package response

type Callback struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}
