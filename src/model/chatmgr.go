package model

import (
	"community-sig/constants"
	"strconv"
	"sync"
	"time"
)

type ChatRoom struct {
	chatRoomId string
	caller     string
	answer     string
}

func (cm *ChatRoom) ID() string {
	return cm.chatRoomId
}

func (cm *ChatRoom) Caller() string {
	return cm.caller
}

func (cm *ChatRoom) Answer() string {
	return cm.answer
}

type ChatMgr struct {
	chatRoomData map[string]*ChatRoom
}

var chatMgrInst *ChatMgr
var chatMgrOnce sync.Once

func GetChatMgr() *ChatMgr {
	chatMgrOnce.Do(func() {
		chatMgrInst = &ChatMgr{
			chatRoomData: make(map[string]*ChatRoom),
		}
	})
	return chatMgrInst
}

//根据chatroom id获取聊天房间信息

func (cm *ChatMgr) GetChatRoom(roomid string) *ChatRoom {
	roomdata, ok := cm.chatRoomData[roomid]
	if !ok {
		return nil
	}
	return roomdata
}

//1对多模式下，创建聊天房间，房间只有两个人

func (cm *ChatMgr) Mul2ChatRoom(roomId_ string, caller_ string, answer_ string) *ChatRoom {
	newRoom := &ChatRoom{caller: caller_, answer: answer_, chatRoomId: roomId_}
	cm.chatRoomData[roomId_] = newRoom
	return newRoom
}

//1对1模式下，创建聊天房间，房间里两个人
func (cm *ChatMgr) CreateChatRoom(caller_ string, answer_ string) *ChatRoom {
	curtime := time.Now().UnixNano() / 1e6
	times := strconv.FormatInt(curtime, 10)
	roomId := caller_ + "-" + answer_ + "-" + times
	newRoom := &ChatRoom{caller: caller_, answer: answer_, chatRoomId: roomId}
	cm.chatRoomData[roomId] = newRoom
	return newRoom
}

//获取聊天对象，根据一个id找到另一个
func (cm *ChatMgr) GetChatMate(roomId string, selfId string) (error, string) {
	roomdata, ok := cm.chatRoomData[roomId]
	if !ok {
		return constants.ErrChatRoomInvalid, ""
	}

	if roomdata.caller == selfId {
		return nil, roomdata.answer
	}

	if roomdata.answer == selfId {
		return nil, roomdata.caller
	}

	return constants.ErrUserInvalid, ""
}

//删除聊天房间里的人员信息
func (cm *ChatMgr) DelChatRoom(roomid string) {
	roomdata, ok := cm.chatRoomData[roomid]
	if !ok {
		return
	}

	answer := GetUserMgr().GetUser(roomdata.Answer())
	if answer != nil {
		answer.SetState(constants.User_Idle)
		answer.SetChatRoomId("")
	}

	caller := GetUserMgr().GetUser(roomdata.Caller())
	if caller != nil {
		caller.SetState(constants.User_Idle)
		caller.SetChatRoomId("")
	}

	delete(cm.chatRoomData, roomid)
}
