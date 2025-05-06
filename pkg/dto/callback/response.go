package callback

type Response struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}
