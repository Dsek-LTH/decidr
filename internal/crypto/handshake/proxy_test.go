package handshake

import (
	"context"
	"testing"
	"time"

	"github.com/flynn/noise"
)

func TestE2EHandshakeThroughProxy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	router := NewRouter()
	clientIdentity, adminIdentity := getIdentityPair(t)

	// 1. Setup Admin Peer
	// adminSide is what the Admin logic uses. proxySideAdmin is what the Proxy uses.
	adminSide, proxySideAdmin := newInMemoryPeers()
	router.RegisterAdmin("admin-1", proxySideAdmin)

	// 2. Setup Client Peer
	clientSide, proxySideClient := newInMemoryPeers()
	router.RegisterClient("client-1", proxySideClient)

	// 3. Start Proxy Forwarding Loops
	// These loops represent what the Web Servers will do:
	// Listen for incoming bytes and call Route.
	go func() {
		for {
			msg, err := proxySideAdmin.Receive(ctx)
			if err != nil {
				return
			}
			// Logic: Admin must prefix message with ClientID or use a protocol
			// for this test, we just route to our one client.
			_ = router.RouteToClient(ctx, "client-1", msg)
		}
	}()

	go func() {
		for {
			msg, err := proxySideClient.Receive(ctx)
			if err != nil {
				return
			}
			_ = router.RouteToAdmin(ctx, "admin-1", msg)
		}
	}()

	// 4. Perform Handshake (The Actual Logic)
	// Note: We use the *side* peers. The handshake happens between Admin and Client.
	clientResCh := runHandshakeAsync(ctx, clientSide, clientIdentity)
	adminResCh := runHandshakeAsync(ctx, adminSide, adminIdentity)

	clientRes := <-clientResCh
	adminRes := <-adminResCh

	if clientRes.err != nil {
		t.Fatalf("Client handshake failed: %v", clientRes.err)
	}
	if adminRes.err != nil {
		t.Fatalf("Admin handshake failed: %v", adminRes.err)
	}

	// 5. Verification
	if string(clientRes.hash) != string(adminRes.hash) {
		t.Error("Handshake hashes do not match through proxy")
	}
}

func TestMultipleAdminsHandshake(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	router := NewRouter()

	// Helper to setup a routed pair
	setupPair := func(adminID, clientID string) (Peer, Peer) {
		clientSide, proxySideClient := newInMemoryPeers()
		adminSide, proxySideAdmin := newInMemoryPeers()

		router.RegisterAdmin(adminID, proxySideAdmin)
		router.RegisterClient(clientID, proxySideClient)

		// Forwarding loops (simulating the web server logic)
		go func() {
			for {
				msg, err := proxySideAdmin.Receive(ctx)
				if err != nil {
					return
				}
				_ = router.RouteToClient(ctx, clientID, msg)
			}
		}()
		go func() {
			for {
				msg, err := proxySideClient.Receive(ctx)
				if err != nil {
					return
				}
				_ = router.RouteToAdmin(ctx, adminID, msg)
			}
		}()

		return adminSide, clientSide
	}

	// Setup two distinct admin-client pairs
	admin1ID, client1ID := "admin-1", "client-1"
	admin2ID, client2ID := "admin-2", "client-2"

	a1Side, c1Side := setupPair(admin1ID, client1ID)
	a2Side, c2Side := setupPair(admin2ID, client2ID)

	// Identities
	c1Identity, a1Identity := getIdentityPair(t)
	c2Identity, a2Identity := getIdentityPair(t)

	// Perform all handshakes concurrently
	resultA1 := runHandshakeAsync(ctx, a1Side, a1Identity)
	resultC1 := runHandshakeAsync(ctx, c1Side, c1Identity)
	resultA2 := runHandshakeAsync(ctx, a2Side, a2Identity)
	resultC2 := runHandshakeAsync(ctx, c2Side, c2Identity)

	// Verify Pair 1
	outA1, outC1 := <-resultA1, <-resultC1
	if outA1.err != nil || outC1.err != nil {
		t.Fatalf("Pair 1 failed: adminErr=%v, clientErr=%v", outA1.err, outC1.err)
	}

	// Verify Pair 2
	outA2, outC2 := <-resultA2, <-resultC2
	if outA2.err != nil || outC2.err != nil {
		t.Fatalf("Pair 2 failed: adminErr=%v, clientErr=%v", outA2.err, outC2.err)
	}

	// Check isolation (Hashes should be unique to the keypairs used)
	if string(outA1.hash) == string(outA2.hash) {
		t.Error("Different admin pairs produced identical handshake hashes")
	}
}

func TestProxyPostHandshakeCommunication(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	router := NewRouter()
	clientIdentity, adminIdentity := getIdentityPair(t)

	// 1. Setup Peers
	adminSide, proxySideAdmin := newInMemoryPeers()
	clientSide, proxySideClient := newInMemoryPeers()

	router.RegisterAdmin("admin-1", proxySideAdmin)
	router.RegisterClient("client-1", proxySideClient)

	// 2. Start Routing Loops
	go func() {
		for {
			msg, err := proxySideAdmin.Receive(ctx)
			if err != nil {
				return
			}
			_ = router.RouteToClient(ctx, "client-1", msg)
		}
	}()
	go func() {
		for {
			msg, err := proxySideClient.Receive(ctx)
			if err != nil {
				return
			}
			_ = router.RouteToAdmin(ctx, "admin-1", msg)
		}
	}()

	// 3. Perform Handshake and capture CipherStates
	type result struct {
		send *noise.CipherState
		recv *noise.CipherState
		err  error
	}

	clientDone := make(chan result, 1)
	adminDone := make(chan result, 1)

	go func() {
		send := func(b []byte) error { return clientSide.Send(ctx, b) }
		receive := func() ([]byte, error) { return clientSide.Receive(ctx) }
		s, r, _, err := Perform(ctx, send, receive, clientIdentity)
		clientDone <- result{s, r, err}
	}()

	go func() {
		send := func(b []byte) error { return adminSide.Send(ctx, b) }
		receive := func() ([]byte, error) { return adminSide.Receive(ctx) }
		s, r, _, err := Perform(ctx, send, receive, adminIdentity)
		adminDone <- result{s, r, err}
	}()

	cRes := <-clientDone
	aRes := <-adminDone

	if cRes.err != nil || aRes.err != nil {
		t.Fatalf("Handshake failed: client=%v, admin=%v", cRes.err, aRes.err)
	}

	// 4. Test Communication: Client -> Admin
	t.Run("ClientToAdmin", func(t *testing.T) {
		plaintext := []byte("secret vote cast")
		ciphertext, err := cRes.send.Encrypt(nil, nil, plaintext)
		if err != nil {
			t.Fatalf("Client failed to encrypt: %v", err)
		}

		// Send through proxy
		go func() { _ = clientSide.Send(ctx, ciphertext) }()

		// Admin receives
		receivedCiphertext, err := adminSide.Receive(ctx)
		if err != nil {
			t.Fatalf("Admin failed to receive: %v", err)
		}

		decrypted, err := aRes.recv.Decrypt(nil, nil, receivedCiphertext)
		if err != nil {
			t.Fatalf("Admin failed to decrypt: %v", err)
		}

		if string(decrypted) != string(plaintext) {
			t.Errorf("Decryption mismatch. Got %s, want %s", decrypted, plaintext)
		}
	})

	// 5. Test Communication: Admin -> Client
	t.Run("AdminToClient", func(t *testing.T) {
		plaintext := []byte("confirmation: received")
		ciphertext, err := aRes.send.Encrypt(nil, nil, plaintext)
		if err != nil {
			t.Fatalf("Admin failed to encrypt: %v", err)
		}

		// Send through transport
		go func() { _ = adminSide.Send(ctx, ciphertext) }()

		// Client receives
		receivedCiphertext, err := clientSide.Receive(ctx)
		if err != nil {
			t.Fatalf("Client failed to receive: %v", err)
		}

		decrypted, err := cRes.recv.Decrypt(nil, nil, receivedCiphertext)
		if err != nil {
			t.Fatalf("Client failed to decrypt: %v", err)
		}

		if string(decrypted) != string(plaintext) {
			t.Errorf("Decryption mismatch. Got %s, want %s", decrypted, plaintext)
		}
	})
}
