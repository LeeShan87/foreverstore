package main

import (
	"log"

	"github.com/leeshan87/foreverstore/p2p"
)

func main() {
	server := p2p.NewTCPTransport(":4000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}
	select {}
}
