package handshake

import "github.com/flynn/noise"

type handshakeIdentity interface {
	getRole() role
	getNoiseConfig() noise.Config
}

type clientIdentity struct {
	AdminPublicKey []byte
}

func (clientIdentity) getRole() role { return initiator }

type adminIdentity struct {
	PublicKey  []byte
	PrivateKey []byte
}

func (adminIdentity) getRole() role { return responder }
