package handshake

import "github.com/Dsek-LTH/decidr/internal/crypto"

type endpoint struct {
	Identity handshakeIdentity
}

func NewAdminEndpoint() (endpoint, endpoint, error) {
	clientID, adminID, err := newAdminClientPair()
	if err != nil {
		return endpoint{}, endpoint{}, err
	}

	return endpoint{Identity: clientID}, endpoint{Identity: adminID}, nil
}

func newAdminClientPair() (clientIdentity, adminIdentity, error) {
	pub, priv, err := crypto.GenerateStaticKeypair()
	if err != nil {
		return clientIdentity{}, adminIdentity{}, err
	}

	return clientIdentity{
			AdminPublicKey: pub,
		},
		adminIdentity{
			PublicKey:  pub,
			PrivateKey: priv,
		},
		nil
}
