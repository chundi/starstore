package message

import (
	"bytes"
	"encoding/json"
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
	//handling *ChMsg
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

func (c *Client) ack(msgId int64, errStr string, ok bool) {
	var res string
	if ok {
		res = ACK_OK
	} else {
		res = ACK_ERROR
	}
	a := Message{
		Id:       msgId,
		Sender:   "",
		Receiver: "",
		Type:     MSG_TYPE_ACK,
		Body: MsgBodyAck{
			Result:  res,
			Message: errStr,
		},
	}
	r, err := json.Marshal(&a)
	if err != nil {
		logger.WithField("storeId", c.owner.id).
			WithField("sender", c.id).
			Error("GENERATE ACK JSON ERROR!!", err)
		return
	}
	c.send <- r
}

func (c *Client) readPump() {
	defer func() {
		c.owner.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(d string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		logger.Info("Pong ", d)
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			}
			logger.Info(err)
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newLine, charSpace, -1))
		msgStr := string(msg)
		logger.WithField("sender", c.id).Info(msgStr)
		chMsg := &ChMsg{
			SenderId: c.id,
			Data:     msg,
			DataStr:  msgStr,
		}
		//c.ack(msgId, "", true)
		//c.handling = chMsg
		c.owner.transfer <- chMsg
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
			logger.Info("Ping ", c.id)
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte(c.id)); err != nil {
				return
			}
		}
	}
}
