package runner

type TransportMsg struct {
	Body     []byte
	Headers  map[string]string
	MetaData map[string]interface{}

	Subject string
}

func (t *TransportMsg) Reply(*TransportMsg) {

}
