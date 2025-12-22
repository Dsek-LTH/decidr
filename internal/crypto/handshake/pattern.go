package handshake

type pattern struct {
	initiator []step
	responder []step
}

var patternNK = pattern{
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
	case initiator:
		return patternNK.initiator
	case responder:
		return patternNK.responder
	default:
		panic("unknown handshake role")
	}
}
