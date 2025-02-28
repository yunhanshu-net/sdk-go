package runner

// Transport 通信层协议
type Transport interface {
	ReadMessage() (*TransportMsg, error)
	WriteMessage(*TransportMsg) error
	Done() <-chan struct{}
	Ping() error
	Close() error
}
