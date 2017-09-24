package message

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type ChMsg struct {
	senderId string
	data     []byte
	dataStr  string
}

type Message struct {
	Id   int64       `json:"id"`
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type MsgBodyExchange struct {
	MainId   string `json:"main_id"`
	SkuId    string `json:"sku_id"`
	DeviceId string `json:"device_id"`
	Space    string `json:"space"`
}

type MsgBodySpace struct {
	DeviceId string `json:"device_id"`
	Space    string `json:"space"`
	Status   string `json:"status"`
}

type MsgBodyAck struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func ProcessMessage(s *Store, m *ChMsg) {
	if !gjson.Valid(m.dataStr) {
		ServerAck(s, m, "Not a valid json.", false)
		return
	}
	switch gjson.Get(m.dataStr, "type").Str {
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
		WithField("sender", m.senderId).
		WithField("data", m.dataStr).
		WithField("info", errStr)
	if !ok {
		flog.Error()
	}
	sender, exist := s.getClient(m.senderId)
	if !exist {
		flog.Error("SEND ACK ERROR, SENDER OFFLINE!!")
		return
	}
	msgId := gjson.Get(m.dataStr, "id").Int()
	sender.ack(msgId, errStr, ok)
}

func ProcessReqExchange(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.senderId).
		WithField("data", m.dataStr)
	sender, exist := s.getClient(m.senderId)
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
		Id:   gjson.Get(m.dataStr, "id").Int(),
		Type: MSG_TYPE_REQ_EXCHANGE,
		Body: MsgBodyExchange{
			MainId:   gjson.Get(m.dataStr, "body.main_id").Str,
			SkuId:    gjson.Get(m.dataStr, "body.sku_id").Str,
			DeviceId: m.senderId,
			Space:    sender.name,
		},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE REQ_EXCHANGE JSON ERROR!!", err)
		return
	}
	receiver.send <- r
	//ServerAck(s, m, "", true)
}

func ProcessRspExchange(s *Store, m *ChMsg) {
	TransferMsg(s, m)
}

func ProcessLsSpace(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.senderId).
		WithField("data", m.dataStr)
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
			DeviceId: client.id,
			Space:    client.name,
			Status:   status,
		}
		spaces = append(spaces, space)
	}
	s.rwLock.RUnlock()
	msg := Message{
		Id:   gjson.Get(m.dataStr, "id").Int(),
		Type: MSG_TYPE_ACK,
		Body: spaces,
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		flog.Error("GENERATE LS_SPACE JSON ERROR!!", err)
		return
	}
	sender, exist := s.getClient(m.senderId)
	if !exist {
		flog.Error("SENDER OFFLINE!!")
		return
	}
	sender.send <- r
}

func ProcessBindSpace(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.senderId).
		WithField("data", m.dataStr)
	sender, exist := s.getClient(m.senderId)
	if !exist {
		flog.Error("SENDER OFFLINE!!")
		return
	}
	if !gjson.Get(m.dataStr, "body.spaces").IsArray() {
		flog.Error("NOT A VALID JSON")
		ServerAck(s, m, "Not a valid json.", false)
		return
	}
	clientIds := gjson.Get(m.dataStr, "body.spaces").Array()
	for _, clientId := range clientIds {
		client, exist := s.getClient(clientId.Str)
		if !exist {
			flog.Error("Bound failed, device not found!", clientId.Str)
			continue
		}
		client.watcher = sender
		sender.watching[clientId.Str] = client
	}
	ServerAck(s, m, "", true)
}

func ProcessTestMsg(s *Store, m *ChMsg) {
	receiverId := gjson.Get(m.dataStr, "body.device_id").Str
	receiver, exist := s.getClient(receiverId)
	if !exist {
		ServerAck(s, m, "Device not found", false)
		logger.Error("TEST MESSAGE, RECEIVER NOT FOUND!!")
		return
	}
	receiver.send <- m.data
	//ServerAck(s, m, "", true)
}

func ProcessCheckIn(s *Store, m *ChMsg) {
	TransferMsg(s, m)
}

func ProcessCheckOut(s *Store, m *ChMsg) {
	TransferMsg(s, m)
}

func ProcessAck(s *Store, m *ChMsg) {
	TransferMsg(s, m)
}

func TransferMsg(s *Store, m *ChMsg) {
	flog := logger.WithField("storeId", s.id).
		WithField("sender", m.senderId).
		WithField("data", m.dataStr)
	receiverId := gjson.Get(m.dataStr, "body.device_id").Str
	receiver, exist := s.getClient(receiverId)
	if !exist {
		flog.Error("RECEIVE DEVICE OFFLINE!!")
		ServerAck(s, m, "Device not found", false)
		return
	}
	receiver.send <- m.data
}
