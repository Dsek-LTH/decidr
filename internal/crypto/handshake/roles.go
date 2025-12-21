package handshake

type role int

func (r role) valid() bool {
	return r == Initiator || r == Responder
}

const (
	Initiator role = iota
	Responder
)
