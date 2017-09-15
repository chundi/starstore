package message

import "time"

const (
	writeWait       = 10 * time.Second
	maxMsgSize      = 1024
	maxReadDeadLine = time.Second * 3
	pongWait        = time.Second * 5
	pingPeriod      = (pongWait * 9) / 10

	MSG_REQ_EXCHANGE = "req_exchange"
	MSG_RSP_EXCHANGE = "rsp_exchange"
	MSG_SERVER_ACK   = "server_ack"
	MSG_CLIENT_ACK   = "client_ack"
	MSG_LS_SPACE     = "ls_space"
	MSG_BIND_SPACE   = "bind_space"
	MSG_TEST         = "test"

	ACK_OK    = "ok"
	ACK_ERROR = "error"
)
