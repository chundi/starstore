package message

import "time"

const (
	writeWait  = 10 * time.Second
	maxMsgSize = 1024
	pongWait   = time.Minute * 3
	pingPeriod = (pongWait * 9) / 10

	MSG_TYPE_REQ_EXCHANGE = "req_exchange"
	MSG_TYPE_RSP_EXCHANGE = "rsp_exchange"
	MSG_TYPE_ACK          = "ack"
	MSG_TYPE_LS_SPACE     = "ls_space"
	MSG_TYPE_RSP_LS_SPACE = "rsp_ls_space"
	MSG_TYPE_LS_REQ       = "ls_req"
	MSG_TYPE_RSP_LS_REQ   = "rsp_ls_req"
	MSG_TYPE_BIND_SPACE   = "bind_space"
	MSG_TYPE_CHECK_IN     = "check_in"
	MSG_TYPE_CHECK_OUT    = "check_out"
	MSG_TYPE_TEST         = "test"

	ACK_OK     = "ok"
	ACK_ERROR  = "error"
	ACK_NOTIFY = "notify"

	SPACE_TYPE_DRESSING_ROOM = "dressing_room"
	SPACE_TYPE_PDA           = "xiaomi"
)
