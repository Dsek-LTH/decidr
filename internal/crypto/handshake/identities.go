package handshake

import "github.com/flynn/noise"

type handshakeIdentity interface {
	GetPublicKey() []byte
	getRole() role
	getNoiseConfig() noise.Config
}

type clientIdentity struct {
	AdminPublicKey []byte
}

func (c clientIdentity) GetPublicKey() []byte { return c.AdminPublicKey }

func (clientIdentity) getRole() role { return initiator }

type adminIdentity struct {
	PublicKey  []byte
	PrivateKey []byte
}

func (a adminIdentity) GetPublicKey() []byte { return a.PublicKey }

func (adminIdentity) getRole() role { return responder }
