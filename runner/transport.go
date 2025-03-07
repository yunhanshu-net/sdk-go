package runner

// Transport 通信层协议
type Transport interface {
	ReadMessage() <-chan *TransportMsg
	Connect() error
	Close() error
}
