package p2p

import "errors"

// ErrIvaligHandshake is returned if the handshake between
// local and remote node could not be enstablished.
var ErrInvalidHandshake = errors.New("invalid handshake")

type HandshakeFunc func(any) error

func NOPHandshake(v any) error {
	return nil
}
