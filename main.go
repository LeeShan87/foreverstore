package main

import (
	"fmt"
	"log"

	"github.com/leeshan87/foreverstore/p2p"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("Handling logic out side TCPTransports")
	return nil
}
func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":4000",
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.NOPHandshake,
		OnPeer:        OnPeer,
	}
	server := p2p.NewTCPTransport(tcpOpts)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}
	go func() {
		for {
			msg := <-server.Consume()
			fmt.Printf("msg: %+v\n", string(msg.Payload))
		}
	}()
	select {}
}
