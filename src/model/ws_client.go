package model

import (
	"community-sig/config"
	"community-sig/logging"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var clientMgrInst *clientMgr
var clientMgrOnce sync.Once

func GetClientMgr() *clientMgr {
	clientMgrOnce.Do(func() {
		clientMgrInst = &clientMgr{}
		clientMgrInst.clientMap = make(map[string]*WSClient)
	})

	return clientMgrInst
}

type WSClient struct {
	id        string
	socket    *websocket.Conn
	closed    bool
	lock      sync.Mutex
	lastHeart int64
}

func (wc *WSClient) SendMsg(data []byte) error {
	logging.Logger.Info("send msg is : ", string(data))
	return wc.socket.WriteMessage(websocket.TextMessage, data)
}

func InitClient(conn *websocket.Conn) *WSClient {
	return &WSClient{id: uuid.NewV4().String(), socket: conn, closed: false, lock: sync.Mutex{}}
}

func (wc *WSClient) UpdateHeartBeat() {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if wc.closed {
		return
	}
	wc.lastHeart = time.Now().Unix()
	logging.Logger.Infof("client %s update heart beat, lastheart is %d", wc.id, wc.lastHeart)
}

func (wc *WSClient) IsAlive(curTime int64) bool {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if (curTime - wc.lastHeart) > int64(config.GetConf().Base.HeartMax) {
		return false
	}

	return true
}

func (wc *WSClient) OnClose() {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if wc.closed {
		return
	}
	wc.socket.Close()
	logging.Logger.Infof("connection %s closed", wc.id)
	wc.closed = true
}

func (wc *WSClient) OnConnect() {
	//增加连接回调
	logging.Logger.Infof("new connection %s connected", wc.id)
	wc.closed = false
	wc.lastHeart = time.Now().Unix()
	//logging.Logger.Info("current goroutine id is ", GetGID())
}

func (wc *WSClient) Id() string {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	return wc.id
}

func (wc *WSClient) IsOnline() bool {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	return !wc.closed
}

type clientMgr struct {
	clientMap map[string]*WSClient
	lock      sync.Mutex
}

func (cm *clientMgr) CheckHeart() []*WSClient {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	var deathIds []*WSClient
	curtime := time.Now().Unix()
	//logging.Logger.Info("begin to check heart")
	for id, client := range cm.clientMap {
		if !client.IsAlive(curtime) {
			logging.Logger.Infof("check %s heart time out, closed the connection", id)
			delete(cm.clientMap, id)
			client.OnClose()
			deathIds = append(deathIds, client)
		}
	}

	return deathIds
}

func (cm *clientMgr) IsClientAlive(client *WSClient) bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	_, ok := cm.clientMap[client.id]
	return ok
}

func (cm *clientMgr) DelClient(client *WSClient) {
	if client == nil {
		return
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()
	client, ok := cm.clientMap[client.id]
	if !ok {
		return
	}
	logging.Logger.Infof("closed client %s", client.id)
	GetC2UMgr().DelC2User(client)
	delete(cm.clientMap, client.id)
	client.OnClose()
}

func (cm *clientMgr) AddClient(client *WSClient) {
	if client == nil {
		return
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.clientMap[client.id] = client
	client.OnConnect()
}
