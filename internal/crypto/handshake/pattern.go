package handshake

type pattern struct {
	initiator []step
	responder []step
}

var patternNN = pattern{
	initiator: []step{
		stepSend{},
		stepReceive{},
	},
	responder: []step{
		stepReceive{},
		stepSend{},
	},
}

func stepsFor(role role) []step {
	switch role {
	case Initiator:
		return patternNN.initiator
	case Responder:
		return patternNN.responder
	default:
		panic("unknown handshake role")
	}
}
