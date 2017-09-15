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
	SkuId    string `json:"sku_id"`
	DeviceId string `json:"device_id"`
	Space    string `json:"space"`
}

type MsgBodySpace struct {
	DeviceId string `json:"device_id"`
	Space    string `json:"space""`
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
	case MSG_REQ_EXCHANGE:
		ProcessReqExchange(s, m)
	case MSG_RSP_EXCHANGE:
		ProcessRspExchange(s, m)
	case MSG_LS_SPACE:
		ProcessLsSpace(s, m)
	case MSG_BIND_SPACE:
		ProcessBindSpace(s, m)
	case MSG_TEST:
		ProcessTestMsg(s, m)
	default:
		ServerAck(s, m, "Unknown message type.", false)
		return
	}
}

func ServerAck(s *Store, m *ChMsg, msg string, ok bool) {
	fields := logger.WithField("storeId", s.id).
		WithField("sender", m.senderId).
		WithField("data", m.dataStr).
		WithField("info", msg)
	if !ok {
		fields.Error()
	}
	sender, exist := s.GetClient(m.senderId)
	if !exist {
		fields.Error("SEND ACK ERROR, SENDER OFFLINE!!")
		return
	}
	var msgId int64 = 0
	if gjson.Valid(m.dataStr) {
		msgId = gjson.Get(m.dataStr, "id").Int()
	}
	var result string
	if ok {
		result = ACK_OK
	} else {
		result = ACK_ERROR
	}
	ack := Message{
		Id:   msgId,
		Type: MSG_SERVER_ACK,
		Body: MsgBodyAck{
			Result:  result,
			Message: msg,
		},
	}
	r, err := json.Marshal(&ack)
	if err != nil {
		fields.Error("GENERATE ACK JSON ERROR!!", err)
		return
	}
	sender.send <- r
}

func ProcessReqExchange(s *Store, m *ChMsg) {
	fields := logger.WithField("storeId", s.id).
		WithField("sender", m.senderId).
		WithField("data", m.dataStr)
	sender, exist := s.GetClient(m.senderId)
	if !exist {
		fields.Error("SENDER OFFLINE!!")
		return
	}
	receiver := sender.watcher
	if receiver == nil {
		fields.Error("NO BIND DEVICES!!")
		ServerAck(s, m, "Receive device not exist!", false)
		return
	}
	msg := Message{
		Id:   gjson.Get(m.dataStr, "id").Int(),
		Type: MSG_REQ_EXCHANGE,
		Body: MsgBodyExchange{
			SkuId:    gjson.Get(m.dataStr, "body.sku_id").Str,
			DeviceId: m.senderId,
			Space:    sender.name,
		},
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		fields.Error("GENERATE REQ_EXCHANGE JSON ERROR!!", err)
		return
	}
	receiver.send <- r
	ServerAck(s, m, "", true)
}

func ProcessRspExchange(s *Store, m *ChMsg) {

}

func ProcessLsSpace(s *Store, m *ChMsg) {
	fields := logger.WithField("storeId", s.id).
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
		Type: MSG_LS_SPACE,
		Body: spaces,
	}
	r, err := json.Marshal(&msg)
	if err != nil {
		fields.Error("GENERATE LS_SPACE JSON ERROR!!", err)
		return
	}
	sender, exist := s.GetClient(m.senderId)
	if !exist {
		fields.Error("SENDER OFFLINE!!")
		return
	}
	sender.send <- r
}

func ProcessBindSpace(s *Store, m *ChMsg) {}

func ProcessTestMsg(s *Store, m *ChMsg) {
	receiverId := gjson.Get(m.dataStr, "body.device_id").Str
	receiver, exist := s.GetClient(receiverId)
	if !exist {
		ServerAck(s, m, "Device not found", false)
		logger.Error("TEST MESSAGE, RECEIVER NOT FOUND!!")
		return
	}
	receiver.send <- m.data
	ServerAck(s, m, "", true)
}
