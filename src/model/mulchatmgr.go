package model

import (
	"community-sig/constants"
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

type MulChatRoom struct {
	chatRoomId  string
	caller      string
	waitAnswers map[string]bool
}

func (cm *MulChatRoom) ID() string {
	return cm.chatRoomId
}

func (cm *MulChatRoom) Caller() string {
	return cm.caller
}

func (cm *MulChatRoom) AnswerValid(anwer string) bool {
	_, ok := cm.waitAnswers[anwer]
	return ok
}

func (cm *MulChatRoom) DeleteWaitAnswer(anwer string) bool {
	delete(cm.waitAnswers, anwer)
	if len(cm.waitAnswers) == 0 {
		return true
	}

	return false
}

type MulChatMgr struct {
	mulChatRoom map[string]*MulChatRoom
}

var mulChatMgrInst *MulChatMgr
var mulChatMgrOnce sync.Once

func GetMulChatMgr() *MulChatMgr {
	mulChatMgrOnce.Do(func() {
		mulChatMgrInst = &MulChatMgr{
			mulChatRoom: make(map[string]*MulChatRoom),
		}
	})
	return mulChatMgrInst
}

func (cm *MulChatMgr) GetMulChatRoom(roomid string) *MulChatRoom {
	mulRoomData, ok := cm.mulChatRoom[roomid]
	if !ok {
		return nil
	}

	return mulRoomData
}

func (cm *MulChatMgr) CreateMulChat(caller_ string, habitRoomId string,
	waiters_ map[string]*UserData) *MulChatRoom {
	curtime := time.Now().UnixNano() / 1e6
	times := strconv.FormatInt(curtime, 10)
	roomId := caller_ + "-" + habitRoomId + "-" + times
	waiterMap := make(map[string]bool)

	for waiter, waiterUd := range waiters_ {
		waiterMap[waiter] = true
		waiterUd.SetState(constants.User_Busy)
		waiterUd.SetChatRoomId(roomId)
		waiterUd.UpdateLastChat()
	}

	callerUd := GetUserMgr().GetUser(caller_)
	callerUd.SetState(constants.User_Busy)
	callerUd.SetChatRoomId(roomId)
	callerUd.UpdateLastChat()

	mulRoom := &MulChatRoom{chatRoomId: roomId, caller: caller_, waitAnswers: waiterMap}

	cm.mulChatRoom[roomId] = mulRoom

	return mulRoom
}

func (cm *MulChatMgr) DelMulRoom(roomid string) {
	mulChatRoom, ok := cm.mulChatRoom[roomid]
	if !ok {
		return
	}

	for answerId, _ := range mulChatRoom.waitAnswers {

		answer := GetUserMgr().GetUser(answerId)
		if answer == nil {
			continue
		}

		answer.SetState(constants.User_Idle)
		answer.SetChatRoomId("")

	}

	caller := GetUserMgr().GetUser(mulChatRoom.caller)
	if caller == nil {
		return
	}

	caller.SetState(constants.User_Idle)
	caller.SetChatRoomId("")

}

func (cm *MulChatMgr) NotifyAnswerHangup(roomid string) {
	mulChatRoom, ok := cm.mulChatRoom[roomid]
	if !ok {
		return
	}

	for answerId, _ := range mulChatRoom.waitAnswers {

		answer := GetUserMgr().GetUser(answerId)
		if answer == nil {
			continue
		}

		answer.SetState(constants.User_Idle)
		answer.SetChatRoomId("")
		if !answer.IsOnline() {
			continue
		}

		hangupNotify := &SCMulHangupNotify{}

		hangupNotifyRsp := ResponseStruct{}
		hangupNotifyRsp.Event = WS_CALL_MULT_HANGUP_NOTIFY
		hangupNotifyRsp.Code = constants.ResponseCodeSuccess
		hangupNotifyRsp.Data = hangupNotify
		hangupNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
		sendata, _ := json.Marshal(hangupNotifyRsp)

		answer.Client.SendMsg(sendata)
	}

	caller := GetUserMgr().GetUser(mulChatRoom.caller)
	if caller == nil {
		return
	}

	caller.SetState(constants.User_Idle)
	caller.SetChatRoomId("")
}

func (cm *MulChatMgr) DelMulRoomWaiters(roomid string, activeAnswer string) {
	mulChatRoom, ok := cm.mulChatRoom[roomid]
	if !ok {
		return
	}

	for answerId, _ := range mulChatRoom.waitAnswers {
		if activeAnswer == answerId {
			continue
		}

		answer := GetUserMgr().GetUser(answerId)
		if answer == nil {
			continue
		}

		answer.SetState(constants.User_Idle)
		answer.SetChatRoomId("")
		if !answer.IsOnline() {
			continue
		}

		otherAccept := &SCMulOtherAcceptNotify{}

		otherAcceptRsp := ResponseStruct{}
		otherAcceptRsp.Event = WS_MULT_OTHER_ACCEPT_NOTIFY
		otherAcceptRsp.Code = constants.ResponseCodeSuccess
		otherAcceptRsp.Data = otherAccept
		otherAcceptRsp.Desc = constants.ResponseCodeSuccess.String()
		sendata, _ := json.Marshal(otherAcceptRsp)

		answer.Client.SendMsg(sendata)
	}

	mulChatRoom.waitAnswers = make(map[string]bool)
	mulChatRoom.waitAnswers[activeAnswer] = true

}
