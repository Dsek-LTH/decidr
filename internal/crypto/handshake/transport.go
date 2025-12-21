package handshake

import "context"

type sender interface {
	Send(context.Context, []byte) error
}

type receiver interface {
	Receive(context.Context) ([]byte, error)
}

type peer interface {
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

func newFuncPeer(
	send func([]byte) error,
	receive func() ([]byte, error),
) peer {
	return funcPeer{
		send:    send,
		receive: receive,
	}
}
