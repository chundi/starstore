package message

import (
	"fmt"
	"sync"
)

type Store struct {
	id         string
	hub        *Hub
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	transfer   chan *ChMsg
	rwLock     sync.RWMutex
}

func newStore(h *Hub, store_id string) *Store {
	return &Store{
		id:         store_id,
		hub:        h,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		transfer:   make(chan *ChMsg),
	}
}

func (s *Store) addClient(client *Client) {
	s.clients[client.id] = client
}

func (s *Store) removeClient(client *Client) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	if cli, ok := s.clients[client.id]; ok {
		if cli.online {
			cli.Reset()
		}
		delete(s.clients, client.id)
	}
}

func (s *Store) getClient(clientId string) (*Client, bool) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	cli, ok := s.clients[clientId]
	return cli, ok
}

func (s *Store) start() {
	for {
		select {
		case client := <-s.register:
			s.addClient(client)
		case client := <-s.unregister:
			info := fmt.Sprintf("who: %s, id: %s offline.", client.name, client.id)
			logger.Info(info)
			if client.online {
				client.Reset()
			}
		case msg := <-s.transfer:
			s.hub.ch_backup <- msg
			go ProcessMessage(s, msg)
		}
	}
}
