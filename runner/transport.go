package runner

// Transport 通信层协议
type Transport interface {
	ReadMessage() <-chan *TransportMsg
	Ping() error
	Close() error
	Wait() <-chan struct{}
}
