// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"fmt"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	seen map[*Client]bool

	// Measurement data to broadcast
	In chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {

	return &Hub{
		//		In:         make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		seen:       make(map[*Client]bool),
	}
}

func (h *Hub) run() {

	finalizeClient := func(client *Client) {
		close(client.send)
		delete(h.seen, client)
		fmt.Println("client removed:", client.conn.RemoteAddr())

	}

	for {
		select {

		case client := <-h.register:

			fmt.Println("client added:", client.conn.RemoteAddr())
			h.seen[client] = true

		case client := <-h.unregister:

			if _, ok := h.seen[client]; ok {
				finalizeClient(client)
			}

		case in := <-h.In:

			for client := range h.seen {

				select {
				case client.send <- in:
				default:
					finalizeClient(client)
				}

			}

		}
	}

}
