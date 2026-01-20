package handshake

import "github.com/Dsek-LTH/decidr/internal/crypto"

type endpoint struct {
	Identity handshakeIdentity
}

func NewAdminEndpoint() (clientEndpoint endpoint, adminEndpoint endpoint, err error) {
	clientID, adminID, err := newAdminClientPair()
	if err != nil {
		return endpoint{}, endpoint{}, err
	}

	return endpoint{Identity: clientID}, endpoint{Identity: adminID}, nil
}

func GetClientEndpoint(adminPublicKey []byte) endpoint {
	return endpoint{
		Identity: clientIdentity{
			AdminPublicKey: adminPublicKey,
		},
	}
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
