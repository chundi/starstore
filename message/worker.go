package message

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/tidwall/gjson"
)

type ChMsg struct {
	Id         int64
	SenderId   string
	ReceiverId string
	Sender     *Client
	Receiver   *Client
	Data       []byte
	DataStr    string
	Msg        *Message
	MsgByte    []byte
}

type Message struct {
	Id       int64       `json:"id"`
	Sender   string      `json:"sender"`
	Receiver string      `json:"receiver"`
	Type     string      `json:"type"`
	Body     interface{} `json:"body"`
}

type MsgBodyReqExchange struct {
	MainId string `json:"main_id"`
	SkuId  string `json:"sku_id"`
	Space  string `json:"space"`
}

type MsgBodyRspExchange struct {
	MainId string `json:"main_id"`
	SkuId  string `json:"sku_id"`
	Status string `json:"status"`
}

type MsgBodySpace struct {
	Receiver string `json:"receiver"`
	Space    string `json:"space"`
	Status   string `json:"status"`
	Watcher  string `json:"watcher"`
}

type MsgBodyAck struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type MsgBodyCheckIn map[string][]string

type MsgBodyEmpty map[string][]string

func PreProcessMessage(s *Store, m *ChMsg) error {
	if !gjson.Valid(m.DataStr) {
		m.Sender.ack(0, "Not a valid json", ACK_ERROR)
		return errors.New("not a valid json!!")
	}
	if gjson.Get(m.DataStr, "id").Exists() {
		m.Id = gjson.Get(m.DataStr, "id").Int()
	}
	//m.Sender.ack(m.Id, "", ACK_OK)
	if gjson.Get(m.DataStr, "receiver").Exists() {
		m.ReceiverId = gjson.Get(m.DataStr, "receiver").Str
	}
	flog := logger.WithField("storeId", s.id).
		WithField("msgId", m.Id).
		WithField("sender", m.SenderId)
	flog.Info(m.DataStr)
	msg := &Message{
		Id:       m.Id,
		Sender:   m.SenderId,
		Receiver: m.ReceiverId,
		Type:     gjson.Get(m.DataStr, "type").Str,
	}
	m.Msg = msg
	return nil
}

func ProcessMessage(s *Store, m *ChMsg) {
	err := PreProcessMessage(s, m)
	if err != nil {
		logger.WithField("data", m.DataStr).Error(err)
		return
	}
	switch m.Msg.Type {
	case MSG_TYPE_REQ_EXCHANGE:
		m.Msg.Body = MsgBodyReqExchange{
			MainId: gjson.Get(m.DataStr, "body.main_id").Str,
			SkuId:  gjson.Get(m.DataStr, "body.sku_id").Str,
			Space:  m.Sender.name,
		}
		ProcessReqExchange(s, m)
	case MSG_TYPE_RSP_EXCHANGE:
		m.Msg.Body = MsgBodyRspExchange{
			MainId: gjson.Get(m.DataStr, "body.main_id").Str,
			SkuId:  gjson.Get(m.DataStr, "body.sku_id").Str,
			Status: gjson.Get(m.DataStr, "body.status").Str,
		}
		ProcessRspExchange(s, m)
	case MSG_TYPE_LS_SPACE:
		ProcessLsSpace(s, m)
	case MSG_TYPE_BIND_SPACE:
		ProcessBindSpace(s, m)
	case MSG_TYPE_TEST:
		ProcessTestMsg(s, m)
	case MSG_TYPE_CHECK_IN:
		ProcessCheckIn(s, m)
	case MSG_TYPE_CHECK_OUT:
		m.Msg.Body = MsgBodyEmpty{}
		ProcessCheckOut(s, m)
	case MSG_TYPE_ACK:
		m.Msg.Body = MsgBodyAck{
			Result:  gjson.Get(m.DataStr, "body.result").Str,
			Message: gjson.Get(m.DataStr, "body.message").Str,
		}
		ProcessAck(s, m)
	case MSG_TYPE_LS_REQ:
		ProcessLsReq(s, m)
	default:
		m.Sender.ack(m.Id, "Unknown message type.", ACK_ERROR)
		return
	}
}

func ProcessReqExchange(s *Store, m *ChMsg) {
	receiver := m.Sender.watcher
	if receiver == nil {
		logger.WithField("msgId", m.Id).Error("DEVICE UNBOUND!!")
		m.Sender.ack(m.Id, "Device unbound!!", ACK_ERROR)
		return
	}
	if !receiver.online {
		logger.WithField("msgId", m.Id).Error("RECEIVER OFFLINE!!")
		m.Sender.ack(m.Id, "Receiver offline!!", ACK_ERROR)
		return
	}
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	m.Sender.handling[m.Msg.Body.(MsgBodyReqExchange).SkuId] = m
	receiver.send <- m
	//ServerAck(s, m, "", true)
}

func ProcessRspExchange(s *Store, m *ChMsg) {
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist || !receiver.online {
		logger.WithField("msgId", m.Id).Error("RECEIVER OFFLINE!!")
		m.Sender.ack(m.Id, "Receiver offline.", ACK_ERROR)
		return
	}
	receiver.send <- m
	delete(receiver.handling, m.Msg.Body.(MsgBodyRspExchange).SkuId)
}

func ProcessLsSpace(s *Store, m *ChMsg) {
	logger.WithField("msgId", m.Id).Error("GEN")
	m.Sender.ack(m.Id, "", ACK_OK)
	s.rwLock.RLock()
	spaces := []MsgBodySpace{}
	for _, client := range s.clients {
		if client.tp != SPACE_TYPE_DRESSING_ROOM {
			continue
		}
		if !client.online {
			continue
		}
		space := MsgBodySpace{
			Receiver: client.id,
			Space:    client.name,
		}
		var status string
		if client.watcher == nil {
			status = "unbound"
			space.Watcher = ""
		} else {
			status = "bound"
			space.Watcher = client.watcher.id
		}
		space.Status = status
		spaces = append(spaces, space)
	}
	s.rwLock.RUnlock()
	m.Msg.Body = spaces
	m.Msg.Sender = ""
	m.Msg.Receiver = m.SenderId
	m.Msg.Type = MSG_TYPE_RSP_LS_SPACE
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	m.Sender.send <- m
}

func ProcessBindSpace(s *Store, m *ChMsg) {
	if !gjson.Get(m.DataStr, "body.spaces").IsArray() {
		logger.WithField("msgId", m.Id).Error("MSG FORMAT ERROR!!")
		m.Sender.ack(m.Id, "Message format error.", ACK_ERROR)
		return
	}
	clientIds := gjson.Get(m.DataStr, "body.spaces").Array()
	for _, clientId := range clientIds {
		client, exist := s.getClient(clientId.Str)
		if !exist {
			logger.WithField("msgId", m.Id).Error("Bound failed, device not found!", clientId.Str)
			m.Sender.ack(m.Id, "Device not found.", ACK_ERROR)
			continue
		}
		client.watcher = m.Sender
		m.Sender.watching[clientId.Str] = client
	}
	m.Sender.ack(m.Id, "", ACK_OK)
}

func ProcessTestMsg(s *Store, m *ChMsg) {
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		logger.WithField("msgId", m.Id).Error("TEST MESSAGE, RECEIVER NOT FOUND!!")
		m.Sender.ack(m.Id, "Device not found", ACK_ERROR)
		return
	}
	if !receiver.online {
		logger.WithField("msgId", m.Id).Error("TEST MESSAGE, RECEIVER OFFLINE!!")
		m.Sender.ack(m.Id, "Receiver offline", ACK_ERROR)
		return
	}
	m.MsgByte = m.Data
	receiver.send <- m
	//ServerAck(s, m, "", true)
}

func ProcessCheckIn(s *Store, m *ChMsg) {
	if !gjson.Get(m.DataStr, "body.sku_ids").IsArray() {
		logger.WithField("msgId", m.Id).Error("MSG FORMAT ERROR!!")
		m.Sender.ack(m.Id, "Msg format error.", ACK_ERROR)
		return
	}
	skuIds := gjson.Get(m.DataStr, "body.sku_ids").Array()
	var skus []string
	for _, skuId := range skuIds {
		skus = append(skus, skuId.Str)
	}
	m.Msg.Body = MsgBodyCheckIn{
		"sku_ids": skus,
	}
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		logger.WithField("msgId", m.Id).Error("RECEIVER OFFLINE!!")
		return
	}
	if !receiver.online {
		logger.WithField("msgId", m.Id).Error("RECEIVER OFFLINE!!")
		m.Sender.ack(m.Id, "Receiver offline", ACK_ERROR)
		return
	}
	receiver.clearHandlingMsg()
	receiver.send <- m
}

func ProcessCheckOut(s *Store, m *ChMsg) {
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist || !receiver.online {
		logger.WithField("msgId", m.Id).Error("RECEIVER OFFLINE!!")
		m.Sender.ack(m.Id, "Receiver offline.", ACK_ERROR)
		return
	}
	receiver.clearHandlingMsg()
	receiver.send <- m
}

func ProcessAck(s *Store, m *ChMsg) {
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		logger.WithField("msgId", m.Id).Error("RECEIVER OFFLINE!!")
		m.Sender.ack(m.Id, "Receiver offline.", ACK_ERROR)
		return
	}
	receiver.send <- m
}

func ProcessLsReq(s *Store, m *ChMsg) {
	m.Sender.ack(m.Id, "", ACK_OK)
	lock := sync.RWMutex{}
	lock.RLock()
	msgs := []*Message{}
	for _, client := range m.Sender.watching {
		for _, msg := range client.handling {
			msgs = append(msgs, msg.Msg)
		}
	}
	lock.RUnlock()
	//m.Msg.Body = msgs
	m.Msg.Body = map[string][]*Message{
		"requests": msgs,
	}
	m.Msg.Sender = ""
	m.Msg.Receiver = m.SenderId
	m.Msg.Type = MSG_TYPE_RSP_LS_REQ
	r, ok := MarshalJson(m.Msg)
	if !ok {
		m.Sender.ack(m.Id, "Server error.", ACK_ERROR)
		return
	}
	m.MsgByte = r
	m.Sender.send <- m
}

func MarshalJson(m *Message) ([]byte, bool) {
	r, err := json.Marshal(m)
	if err != nil {
		logger.WithField("msgId", m.Id).Error("GENERATE JSON ERROR!!")
		return nil, false
	}
	return r, true
}
