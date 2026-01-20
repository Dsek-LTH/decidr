package handshake

import (
	"context"
	"errors"
	"sync"
)

type Router struct {
	admins      map[string]Peer
	adminsMutex sync.RWMutex

	clients      map[string]Peer
	clientsMutex sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		admins:  make(map[string]Peer),
		clients: make(map[string]Peer),
	}
}

// RegisterAdmin sets the peer representing the admin connection.
func (router *Router) RegisterAdmin(id string, peer Peer) {
	router.adminsMutex.Lock()
	defer router.adminsMutex.Unlock()
	router.admins[id] = peer
}

// RemoveAdmin cleans up an admin peer.
func (router *Router) RemoveAdmin(id string) {
	router.adminsMutex.Lock()
	defer router.adminsMutex.Unlock()
	delete(router.admins, id)
}

// RegisterClient adds a client peer to the routing table.
func (router *Router) RegisterClient(id string, peer Peer) {
	router.clientsMutex.Lock()
	defer router.clientsMutex.Unlock()
	router.clients[id] = peer
}

// RemoveClient cleans up a client peer.
func (router *Router) RemoveClient(id string) {
	router.clientsMutex.Lock()
	defer router.clientsMutex.Unlock()
	delete(router.clients, id)
}

// RouteToAdmin takes a message from a client and forwards it to a specific admin.
func (router *Router) RouteToAdmin(ctx context.Context, adminID string, data []byte) error {
	router.adminsMutex.RLock()
	admin, ok := router.admins[adminID]
	router.adminsMutex.RUnlock()

	if !ok {
		return errors.New("target admin not connected")
	}
	return admin.Send(ctx, data)
}

// RouteToClient takes a message from the admin and forwards it to a specific client.
func (router *Router) RouteToClient(ctx context.Context, clientID string, data []byte) error {
	router.clientsMutex.RLock()
	client, ok := router.clients[clientID]
	router.clientsMutex.RUnlock()

	if !ok {
		return errors.New("client not found")
	}
	return client.Send(ctx, data)
}
