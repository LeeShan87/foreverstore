package p2p

// Peer is an interface represent a remote node
type Peer interface {
	Close() error
}

// Transport is anything representing communication
// between nodes in the network. This can be of the
// form (TCP, UDP, websocket, ...)
type Transport interface {
	ListenAndServe() error
	Consume() <-chan RPC
}
