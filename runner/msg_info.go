package runner

const (
	MsgTypeClose          = "MsgClose"
	MsgTypeHeartbeatCheck = "HeartbeatCheck"
	MsgTypeRun            = "Run"
)

type MsgHeader map[string][]string

func (h MsgHeader) Add(key, value string) {
	h[key] = append(h[key], value)
}

// Set sets the header entries associated with key to the single
// element value. It is case-sensitive and replaces any existing
// values associated with key.
func (h MsgHeader) Set(key, value string) {
	h[key] = []string{value}
}

// Get gets the first value associated with the given key.
// It is case-sensitive.
func (h MsgHeader) Get(key string) string {
	if h == nil {
		return ""
	}
	if v := h[key]; v != nil {
		return v[0]
	}
	return ""
}

// Values returns all values associated with the given key.
// It is case-sensitive.
func (h MsgHeader) Values(key string) []string {
	return h[key]
}

// Del deletes the values associated with a key.
// It is case-sensitive.
func (h MsgHeader) Del(key string) {
	delete(h, key)
}
