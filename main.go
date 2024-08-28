package main

import (
	"log"

	"github.com/leeshan87/foreverstore/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":4000",
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.NOPHandshake,
	}
	server := p2p.NewTCPTransport(tcpOpts)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}
	select {}
}
