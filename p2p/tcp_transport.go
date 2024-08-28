package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP enstablished connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrive a connection conn => outbound == true
	// if we accept and retrive a connection conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements the Peer interface
func (tp *TCPPeer) Close() error {
	return tp.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndServe() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.acceptLoop()
	return nil
}

// Consume implements the transport interface, which will return a read-only channel
// for reading the incoming messages recived from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection %s", err)
		conn.Close()
	}()
	fmt.Printf("New TCP conncetion from [%+v]\n", conn)
	if err = t.HandShakeFunc(conn); err != nil {
		return
	}
	peer := NewTCPPeer(conn, false)
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	// read loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP decoder error %v\n", err)
			return
		}

		// fmt.Printf("message: %+v\n", string(rpc.Payload))
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

	//t.peers[conn.RemoteAddr()] = NewTCPPeer(conn, false)

}
