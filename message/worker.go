package message

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type ChMsg struct {
	Id         int64
	SenderId   string
	ReceiverId string
	Data       []byte
	DataStr    string
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
}

type MsgBodyAck struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type MsgBodyCheckIn map[string][]string

type MsgBodyEmpty map[string][]string

func ProcessMessage(s *Store, m *ChMsg) {
	if !gjson.Valid(m.DataStr) {
		ServerAck(s, m, "Not a valid json.", false)
		return
	}
	if gjson.Get(m.DataStr, "receiver").Exists() {
		m.ReceiverId = gjson.Get(m.DataStr, "receiver").Str
	}
	if gjson.Get(m.DataStr, "id").Exists() {
		m.Id = gjson.Get(m.DataStr, "id").Int()
	}
	switch gjson.Get(m.DataStr, "type").Str {
	case MSG_TYPE_REQ_EXCHANGE:
		ProcessReqExchange(s, m)
	case MSG_TYPE_RSP_EXCHANGE:
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
		ProcessCheckOut(s, m)
	case MSG_TYPE_ACK:
		ProcessAck(s, m)
	default:
		ServerAck(s, m, "Unknown message type.", false)
		return
	}
}

func ServerAck(s *Store, m *ChMsg, errStr string, ok bool) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("receiver", m.ReceiverId).
		WithField("Data", m.DataStr).
		WithField("info", errStr)
	if !ok {
		flog.Error()
	}
	sender, exist := s.getClient(m.SenderId)
	if !exist {
		flog.Error("SEND ACK ERROR, SENDER OFFLINE!!")
		return
	}
	sender.ack(m.Id, errStr, ok)
}

func ProcessReqExchange(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("receiver", m.ReceiverId).
		WithField("Data", m.DataStr)
	sender, exist := s.getClient(m.SenderId)
	if !exist {
		flog.Error("SENDER OFFLINE!!")
		return
	}
	receiver := sender.watcher
	if receiver == nil {
		flog.Error("DEVICE UNBOUND!!")
		ServerAck(s, m, "Device unbound!!", false)
		return
	}
	msg := Message{
		Id:       m.Id,
		Sender:   m.SenderId,
		Receiver: m.ReceiverId,
		Type:     MSG_TYPE_REQ_EXCHANGE,
		Body: MsgBodyReqExchange{
			MainId: gjson.Get(m.DataStr, "body.main_id").Str,
			SkuId:  gjson.Get(m.DataStr, "body.sku_id").Str,
			Space:  sender.name,
		},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE REQ_EXCHANGE JSON ERROR!!", err)
		ServerAck(s, m, "Server error.", false)
		return
	}
	receiver.send <- r
	//ServerAck(s, m, "", true)
}

func ProcessRspExchange(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("receiver", m.ReceiverId).
		WithField("data", m.DataStr)
	flog.Info()
	msg := &Message{
		Id:       m.Id,
		Sender:   m.SenderId,
		Receiver: m.ReceiverId,
		Type:     MSG_TYPE_RSP_EXCHANGE,
		Body: MsgBodyRspExchange{
			MainId: gjson.Get(m.DataStr, "body.main_id").Str,
			SkuId:  gjson.Get(m.DataStr, "body.sku_id").Str,
			Status: gjson.Get(m.DataStr, "body.status").Str,
		},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE RSP_EXCHANGE JSON ERROR!!")
		ServerAck(s, m, "Server error.", false)
		return
	}
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		flog.Error("RECEIVER OFFLINE!!")
		ServerAck(s, m, "Receiver offline.", false)
		return
	}
	receiver.send <- r
}

func ProcessLsSpace(s *Store, m *ChMsg) {
	ServerAck(s, m, "", true)
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("Data", m.DataStr)
	s.rwLock.RLock()
	spaces := []MsgBodySpace{}
	for _, client := range s.clients {
		if client.tp != "dressing_room" {
			continue
		}
		var status string
		if client.watcher == nil {
			status = "unbound"
		} else {
			status = "bound"
		}
		space := MsgBodySpace{
			Receiver: client.id,
			Space:    client.name,
			Status:   status,
		}
		spaces = append(spaces, space)
	}
	s.rwLock.RUnlock()
	msg := Message{
		Id:       m.Id,
		Sender:   "",
		Receiver: "",
		Type:     MSG_TYPE_RSP_LS_SPACE,
		Body:     spaces,
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE LS_SPACE JSON ERROR!!", err)
		ServerAck(s, m, "Server error.", false)
		return
	}
	sender, exist := s.getClient(m.SenderId)
	if !exist {
		flog.Error("SENDER OFFLINE!!")
		return
	}
	sender.send <- r
}

func ProcessBindSpace(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("Data", m.DataStr)
	sender, exist := s.getClient(m.SenderId)
	if !exist {
		flog.Error("SENDER OFFLINE!!")
		return
	}
	if !gjson.Get(m.DataStr, "body.spaces").IsArray() {
		flog.Error("MSG FORMAT ERROR!!")
		ServerAck(s, m, "Message format error.", false)
		return
	}
	clientIds := gjson.Get(m.DataStr, "body.spaces").Array()
	for _, clientId := range clientIds {
		client, exist := s.getClient(clientId.Str)
		if !exist {
			flog.Error("Bound failed, device not found!", clientId.Str)
			ServerAck(s, m, "Device not found.", false)
			continue
		}
		client.watcher = sender
		sender.watching[clientId.Str] = client
	}
	ServerAck(s, m, "", true)
}

func ProcessTestMsg(s *Store, m *ChMsg) {
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		ServerAck(s, m, "Device not found", false)
		logger.Error("TEST MESSAGE, RECEIVER NOT FOUND!!")
		return
	}
	receiver.send <- m.Data
	//ServerAck(s, m, "", true)
}

func ProcessCheckIn(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("receiver", m.ReceiverId).
		WithField("data", m.DataStr)
	flog.Info()
	if !gjson.Get(m.DataStr, "body.sku_ids").IsArray() {
		flog.Error("MSG FORMAT ERROR!!")
		ServerAck(s, m, "Msg format error.", false)
		return
	}
	skuIds := gjson.Get(m.DataStr, "body.sku_ids").Array()
	var skus []string
	for _, skuId := range skuIds {
		skus = append(skus, skuId.Str)
	}
	msg := Message{
		Id:       m.Id,
		Sender:   m.SenderId,
		Receiver: m.ReceiverId,
		Type:     MSG_TYPE_CHECK_IN,
		Body: MsgBodyCheckIn{
			"sku_ids": skus,
		},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE LS_SPACE JSON ERROR!!", err)
		return
	}
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		flog.Error("RECEIVER OFFLINE!!")
		return
	}
	receiver.send <- r
}

func ProcessCheckOut(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("receiver", m.ReceiverId).
		WithField("data", m.DataStr)
	flog.Info()
	msg := Message{
		Id:       m.Id,
		Sender:   m.SenderId,
		Receiver: m.ReceiverId,
		Type:     MSG_TYPE_CHECK_OUT,
		Body:     MsgBodyEmpty{},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE JSON ERROR!!")
		ServerAck(s, m, "Server error.", false)
		return
	}
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		flog.Error("RECEIVER OFFLINE!!")
		ServerAck(s, m, "Receiver offline.", false)
		return
	}
	receiver.send <- r
}

func ProcessAck(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.SenderId).
		WithField("receiver", m.ReceiverId).
		WithField("data", m.DataStr)
	flog.Info()
	msg := Message{
		Id:       m.Id,
		Sender:   m.SenderId,
		Receiver: m.ReceiverId,
		Type:     MSG_TYPE_ACK,
		Body: MsgBodyAck{
			Result:  gjson.Get(m.DataStr, "body.result").Str,
			Message: gjson.Get(m.DataStr, "body.message").Str,
		},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE JSON ERROR!!")
		ServerAck(s, m, "Server error.", false)
		return
	}
	receiver, exist := s.getClient(m.ReceiverId)
	if !exist {
		flog.Error("RECEIVER OFFLINE!!")
		ServerAck(s, m, "Receiver offline.", false)
		return
	}
	receiver.send <- r
}
