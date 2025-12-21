package handshake

type sender interface {
	Send([]byte) error
}

type receiver interface {
	Receive() ([]byte, error)
}

type peer interface {
	sender
	receiver
}

type funcPeer struct {
	send    func([]byte) error
	receive func() ([]byte, error)
}

func (p funcPeer) Send(b []byte) error {
	return p.send(b)
}

func (p funcPeer) Receive() ([]byte, error) {
	return p.receive()
}

func newFuncPeer(
	send func([]byte) error,
	receive func() ([]byte, error),
) peer {
	return funcPeer{
		send:    send,
		receive: receive,
	}
}
