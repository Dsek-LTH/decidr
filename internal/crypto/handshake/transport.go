package handshake

import "context"

type sender interface {
	Send(context.Context, []byte) error
}

type receiver interface {
	Receive(context.Context) ([]byte, error)
}

type Peer interface {
	sender
	receiver
}

type funcPeer struct {
	send    func([]byte) error
	receive func() ([]byte, error)
}

func (p funcPeer) Send(ctx context.Context, b []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return p.send(b)
	}
}

func (p funcPeer) Receive(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return p.receive()
	}
}

func NewFuncPeer(
	send func([]byte) error,
	receive func() ([]byte, error),
) Peer {
	return funcPeer{
		send:    send,
		receive: receive,
	}
}
