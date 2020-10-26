package model

import (
	"community-sig/constants"
	"community-sig/logging"
	"sync"
	"time"
)

type UserData struct {
	UserId     string
	UserName   string
	Phone      string
	Avator     string
	RoomList   []string
	State      int
	ChatRoomId string
	Client     *WSClient
	LastChat   int64
}

func (ud *UserData) IsOnline() bool {
	if ud.Client == nil {
		return false
	}

	return ud.Client.IsOnline()
}

func (ud *UserData) OffLine() {
	ud.Client = nil
}

func (ud *UserData) SetState(state int) {
	ud.State = state
}

func (ud *UserData) GetState() int {
	return ud.State
}

func (ud *UserData) SetChatRoomId(roomid string) {
	ud.ChatRoomId = roomid
}

func (ud *UserData) GetChatRoomId() string {
	return ud.ChatRoomId
}

func (ud *UserData) UpdateLastChat() {
	cur := time.Now().Unix()
	ud.LastChat = cur
}

func (ud *UserData) ChatForbid() bool {
	cur := time.Now().Unix()
	if (cur - ud.LastChat) < 3600 {
		return true
	}
	return false
}

type UserMgr struct {
	UserMap map[string]*UserData
}

type RoomData struct {
	UserMap map[string]*UserData
}

type RoomMgr struct {
	RoomMap map[string]*RoomData
}

var um *UserMgr
var umOnce sync.Once

var rm *RoomMgr
var rmOnce sync.Once

var c2um *Client2UserMgr
var c2umOnce sync.Once

var reconnectKick chan *WSClient

func init() {
	reconnectKick = make(chan *WSClient, constants.Recon_Chan_Size)
}

func PutIntoReconKick(oldclient *WSClient) {
	if oldclient == nil {
		return
	}
	select {
	case reconnectKick <- oldclient:
		logging.Logger.Debugf("put client %s into recon kick\n", oldclient.id)
		return
	default:
		select {
		case <-time.After(time.Second):
			logging.Logger.Info("put into recon kick time out ")
			return
		case reconnectKick <- oldclient:
			return
		}
	}
}

func GetFromReconKick() chan *WSClient {
	return reconnectKick
}

func GetUserMgr() *UserMgr {
	umOnce.Do(func() {
		um = &UserMgr{UserMap: make(map[string]*UserData)}
	})
	return um
}

func GetRoomMgr() *RoomMgr {
	rmOnce.Do(func() {
		rm = &RoomMgr{RoomMap: make(map[string]*RoomData)}
	})
	return rm
}

func (um *UserMgr) AddUser(ud *UserData) {
	//判断map中是否存在用户，如果存在需要做踢人操作
	udold, ok := um.UserMap[ud.UserId]
	if ok {
		delete(um.UserMap, udold.UserId)
		GetRoomMgr().DelRoomUser(udold)
		if udold.Client != nil {
			GetC2UMgr().DelC2User(udold.Client)
			//发送踢人chan
			PutIntoReconKick(udold.Client)
		}

	}

	um.UserMap[ud.UserId] = ud
	GetRoomMgr().AddRoom(ud)
	GetC2UMgr().AddC2User(ud.Client, ud)
}

func (um *UserMgr) KickUser(uid string) {
	udold, ok := um.UserMap[uid]
	if !ok {
		return
	}
	delete(um.UserMap, udold.UserId)
	GetRoomMgr().DelRoomUser(udold)
	if udold.Client != nil {
		GetC2UMgr().DelC2User(udold.Client)
		//发送踢人chan
		PutIntoReconKick(udold.Client)
	}
}

func (um *UserMgr) GetUser(uid string) *UserData {
	ud, ok := um.UserMap[uid]
	if !ok {
		return nil
	}
	return ud
}

func (um *RoomMgr) GetRoom(roomId string) *RoomData {
	rm, ok := um.RoomMap[roomId]
	if !ok {
		return nil
	}

	return rm
}

func (rm *RoomMgr) AddRoom(ud *UserData) {
	for _, room := range ud.RoomList {
		rmd, ok := rm.RoomMap[room]
		if !ok {
			rmdt := &RoomData{}
			rmdt.UserMap = make(map[string]*UserData)
			rmdt.UserMap[ud.UserId] = ud
			rm.RoomMap[room] = rmdt
			continue
		}
		rmd.UserMap[ud.UserId] = ud
	}
}

//删除房间内的住户
func (rm *RoomMgr) DelRoomUser(ud *UserData) {
	for _, room := range ud.RoomList {
		roomdata, ok := rm.RoomMap[room]
		if !ok {
			continue
		}

		delete(roomdata.UserMap, ud.UserId)
	}
}

type Client2UserMgr struct {
	c2uMap map[string]*UserData
}

func (cm *Client2UserMgr) AddC2User(client *WSClient, userData *UserData) {
	if client == nil {
		return
	}
	cm.c2uMap[client.id] = userData
}

func (cm *Client2UserMgr) DelC2User(client *WSClient) {
	if client == nil {
		return
	}
	delete(cm.c2uMap, client.Id())
}

func (cm *Client2UserMgr) GetUserByClient(id string) *UserData {
	ud, ok := cm.c2uMap[id]
	if !ok {
		return nil
	}
	return ud
}

func GetC2UMgr() *Client2UserMgr {
	c2umOnce.Do(func() {
		c2um = &Client2UserMgr{c2uMap: make(map[string]*UserData)}
	})
	return c2um
}
