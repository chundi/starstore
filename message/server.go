package message

import (
	"fmt"
	"net/http"

	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/log"
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/model/earth"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	logCfg    = conf.AppConfig.Log
	logger    *logrus.Entry
	newLine   = []byte{'\n'}
	charSpace = []byte{' '}
)

var hub *Hub

func init() {
	hub = newHub()
	go hub.run()
	logger = log.NewLogger(logCfg.Format, logCfg.Level, logCfg.Output).WithField("module", "message")
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		e := fmt.Sprintf("Failed to set websocket upgrade: %+v", err)
		logger.Error(e)
		conn.Close()
		return
	}

	clientId := r.URL.Query().Get("token")
	device := &earth.Device{}
	exist := earth.GetDeviceByToken(model.DB.New(), clientId, device)
	if !exist {
		e := fmt.Sprintf("Device %s Not Found!", clientId)
		logger.Error(e)
		conn.WriteMessage(websocket.CloseMessage, []byte(e))
		conn.Close()
		return
	}
	storeId := device.OwnerId
	clientName := device.Name
	clientType := device.Type
	logger.Info(r.RequestURI)

	store, ok := hub.GetStore(storeId)
	if !ok {
		store = newStore(hub, storeId)
		go store.start()
		hub.store_add <- store
	}
	newCli := newClient(clientId, store, conn, clientName, clientType)
	client, ok := store.getClient(clientId)
	if ok {
		logger.Info(fmt.Sprintf("Client %s %s ReConnecting, close old connection.", clientName, clientId))
		newCli.watching = client.watching
		newCli.handling = client.handling
		newCli.watcher = client.watcher
		for _, watchingId := range client.watching {
			cli, exist := store.getClient(watchingId)
			if !exist {
				continue
			}
			cli.watcher = clientId
		}
		if client.online {
			client.Reset()
		}
		client.watcher = ""
		client.handling = nil
		client.watching = nil
	}
	store.addClient(newCli)
	go newCli.readPump()
	go newCli.writePump()
}
