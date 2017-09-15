package message

import (
	"bytes"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	id       string
	tp       string
	owner    *Store
	name     string
	watcher  *Client
	watching map[string]*Client
	conn     *websocket.Conn
	send     chan []byte
}

func newClient(id string, store *Store, conn *websocket.Conn, name string, tp string) *Client {
	return &Client{
		id:       id,
		owner:    store,
		tp:       tp,
		name:     name,
		conn:     conn,
		watching: make(map[string]*Client),
		send:     make(chan []byte),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.owner.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newLine, charSpace, -1))
		logger.WithField("sender", c.id).Info(string(msg))
		c.owner.transfer <- &ChMsg{c.id, msg, string(msg)}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
				w.Write(newLine)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
