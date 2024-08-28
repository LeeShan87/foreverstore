package p2p

type HandshakeFunc func(any) error

func NOPHandshake(v any) error {
	return nil
}
