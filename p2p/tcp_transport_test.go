package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOpt := TCPTransportOpts{
		ListenAddr:    ":4000",
		Decoder:       DefaultDecoder{},
		HandShakeFunc: NOPHandshake,
	}
	listenAddr := ":4000"
	tr := NewTCPTransport(tcpOpt)
	assert.Equal(t, listenAddr, tcpOpt.ListenAddr)
	assert.Nil(t, tr.ListenAndServe())
}
