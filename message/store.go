package message

import (
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
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	if cli, ok := s.clients[client.id]; ok {
		cli.conn.Close()
		close(cli.send)
		delete(s.clients, client.id)
	}
	s.clients[client.id] = client
}

func (s *Store) removeClient(client *Client) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	if cli, ok := s.clients[client.id]; ok {
		close(cli.send)
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
			s.removeClient(client)
		case msg := <-s.transfer:
			go ProcessMessage(s, msg)
		}
	}
}
