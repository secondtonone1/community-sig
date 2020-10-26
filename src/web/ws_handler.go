package web

import (
	"community-sig/constants"
	"community-sig/grpc_client"
	"community-sig/logging"
	"community-sig/model"
	"community-sig/protobuffer_def"
	"encoding/json"
	"sync"

	"community-sig/service/impl"

	"github.com/goinggo/mapstructure"
	"github.com/golang/protobuf/ptypes"
)

var webMsgMap map[string]handlerWraper
var weblogicLock sync.Mutex
var heartExit chan struct{}
var closeHeart chan struct{}

func init() {
	webMsgMap = make(map[string]handlerWraper)
	heartExit = make(chan struct{})
	closeHeart = make(chan struct{})
	RegMsgHandler(model.WS_Login_CS, login_CS)
	RegMsgHandler(model.WS_HEART_BEAT_CS, heart_Beat)
	RegMsgHandler(model.WS_CALL_SINGLE_CS, single_Call)
	RegMsgHandler(model.WS_CALL_SINGLE_ANSWER_CS, single_Answer)
	RegMsgHandler(model.WS_CALL_SINGLE_REFUSE_CS, single_Refuse)
	RegMsgHandler(model.WS_SINGLE_HANGUP_CS, single_Hangup)
	RegMsgHandler(model.WS_SINGLE_TERMINATE_CS, single_Terminal)
	RegMsgHandler(model.WS_Offline_SYS, offine_SYS)
	RegMsgHandler(model.WS_OFFER_CALL_CS, offer_Call)
	RegMsgHandler(model.WS_OFFER_ANSWER_CS, offer_Answer)
	RegMsgHandler(model.WS_OFFER_CALL_CS, offer_Call)
	RegMsgHandler(model.WS_OFFER_ANSWER_CS, offer_Answer)
	RegMsgHandler(model.WS_ICE_CALL_CS, ice_Call)
	RegMsgHandler(model.WS_ICE_ANSWER_CS, ice_Answer)
	RegMsgHandler(model.WS_MEDIA_TO_AUDIO_CS, media_ToAudio)
	RegMsgHandler(model.WS_CALL_MULT_CS, call_Mul)
	RegMsgHandler(model.WS_CALL_MULT_ANSWER_CS, call_MulAnswer)
	RegMsgHandler(model.WS_CALL_MULT_REFUSE_CS, call_MulRefuse)
	RegMsgHandler(model.WS_CALL_MULT_HANGUP_CS, call_Hangup)
}

type handlerWraper interface {
	HandleMsg(*model.WSClient, interface{}) error
}

type handlerWraperImpl struct {
	handler func(*model.WSClient, interface{}) error
	event   string
}

func (hw *handlerWraperImpl) HandleMsg(ws *model.WSClient, data interface{}) error {
	weblogicLock.Lock()
	defer weblogicLock.Unlock()
	return hw.handler(ws, data)
}

func RegMsgHandler(event_ string, handler_ func(*model.WSClient, interface{}) error) {
	hw := &handlerWraperImpl{event: event_, handler: handler_}
	webMsgMap[event_] = hw
}

func IfMsgReged(event string) bool {
	if _, ok := webMsgMap[event]; !ok {
		return false
	}
	return true
}

func GetMsgHandler(event string) handlerWraper {
	handler, ok := webMsgMap[event]
	if !ok {
		return nil
	}

	return handler
}

func login_CS(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive login_cs req")
	loginReq := &model.CSLogin{}
	if err := mapstructure.Decode(data, loginReq); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err.Error())
		return constants.ErrMap2Struct
	}

	logging.Logger.Info("loginReq.UserId is", loginReq.UserId)
	logging.Logger.Info("loginReq.RoomList is", loginReq.RoomList)
	loginRsp := &model.ResponseStruct{}
	impl.GetSigServiceImpl().UserLogin(client, loginReq, loginRsp)

	sendata, _ := json.Marshal(loginRsp)

	return client.SendMsg(sendata)
}

func heart_Beat(client *model.WSClient, data interface{}) error {
	//logging.Logger.Info("receive heart_Beat req")
	heartReq := &model.CSHeartBeat{}
	if err := mapstructure.Decode(data, heartReq); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	//logging.Logger.Info("heartReq.UserId is ", heartReq.UserId)
	return impl.GetSigServiceImpl().UpdateHeartBeat(client, heartReq)
}

func single_Call(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive single_Call req")
	callSingle := &model.CSCallSingle{}
	if err := mapstructure.Decode(data, callSingle); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}
	callSingleRsp := &model.ResponseStruct{}
	err := impl.GetSigServiceImpl().CallSingle(client, callSingle, callSingleRsp)
	if err != nil {
		return constants.ErrCallSingle
	}
	sendata, err := json.Marshal(callSingleRsp)
	if err != nil {
		return constants.ErrJsonMarshal
	}

	err = client.SendMsg(sendata)
	if err != nil {
		return constants.ErrSendData
	}

	return nil

}

//向status服务器同步用户下线信息
func RPCUpdateOffLine(userId string) error {
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "update_off_line"
	baseReq.C = protobuffer_def.CMD_UPDATE_OFF_LINE
	queryReq := &protobuffer_def.UpdateOfflineReq{}
	queryReq.Identity = userId

	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return err
	}
	baseReq.Body = body

	logging.Logger.Info("update off line basereq  is ", baseReq)
	_, err = grpc_client.SendRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("update off line failed, err is ", err)
		return err
	}

	return nil
}

func offine_SYS(client *model.WSClient, data interface{}) error {
	ud := model.GetC2UMgr().GetUserByClient(client.Id())
	if ud == nil {
		logging.Logger.Info("client to user nil,  client id is ", client.Id())
		return nil
	}

	if ud.Client != client {
		logging.Logger.Info("client param is ", client)
		logging.Logger.Info("ud.client is ", ud.Client)
		logging.Logger.Infof("user %s has been reconnect, client has chanded\n", ud.UserId)
		return nil
	}
	//设置用户下线
	ud.OffLine()
	//更新下线信息到status服务
	RPCUpdateOffLine(ud.UserId)
	/*

		if ud.ChatRoomId == "" {
			logging.Logger.Infof("user %s state set offline ", ud.UserId)
			return nil
		}

		defer func() {
			model.GetChatMgr().DelChatRoom(ud.ChatRoomId)
			model.GetMulChatMgr().DelMulRoom(ud.ChatRoomId)
			logging.Logger.Infof("user %s state set offline ", ud.UserId)
		}()
		//清理聊天室
		err, mateId := model.GetChatMgr().GetChatMate(ud.ChatRoomId, ud.UserId)
		if err != nil {
			return nil
		}

		mater := model.GetUserMgr().GetUser(mateId)
		if mater == nil {
			return nil
		}

		if mater.IsOnline() == false {
			return nil
		}

		terminalNotify := &model.SCForceTerminateNotify{ChatRoomId: ud.ChatRoomId}

		terminalRsp := model.ResponseStruct{}
		terminalRsp.Event = model.WS_FORCE_TERMINATE_NOTIFY
		terminalRsp.Code = constants.ResponseCodeSuccess
		terminalRsp.Data = terminalNotify
		terminalRsp.Desc = constants.ResponseCodeSuccess.String()
		sendata, _ := json.Marshal(terminalRsp)

		mater.Client.SendMsg(sendata)
	*/
	return nil
}

func single_Answer(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive single_Answer req")
	singleAnswer := &model.CSAnswerSingle{}

	if err := mapstructure.Decode(data, singleAnswer); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	singleAnswerRsp := &model.ResponseStruct{}

	return impl.GetSigServiceImpl().AnswerSingle(client, singleAnswer, singleAnswerRsp)

}

func single_Refuse(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive single_Refuse req")
	singleRefuse := &model.CSRefuseSingle{}

	if err := mapstructure.Decode(data, singleRefuse); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	singleRefuseRsp := &model.ResponseStruct{}

	return impl.GetSigServiceImpl().RefuseSingle(client, singleRefuse, singleRefuseRsp)

}

func single_Hangup(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive single_Hangup req")
	singleHangup := &model.CSHangupSingle{}

	if err := mapstructure.Decode(data, singleHangup); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	singleHangupRsp := &model.ResponseStruct{}

	return impl.GetSigServiceImpl().HangupSingle(client, singleHangup, singleHangupRsp)

}

func single_Terminal(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive single_Terminal req")
	singleTerminal := &model.CSTerminateSingle{}

	if err := mapstructure.Decode(data, singleTerminal); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	return impl.GetSigServiceImpl().TerminalSingle(client, singleTerminal)
}

func offer_Call(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive offer_Call req")
	offerCall := &model.CSOfferCall{}
	if err := mapstructure.Decode(data, offerCall); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	return impl.GetSigServiceImpl().OfferCall(client, offerCall)
}

func offer_Answer(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive offer_Answer req")
	offerAnswer := &model.CSOfferAnswer{}
	if err := mapstructure.Decode(data, offerAnswer); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	return impl.GetSigServiceImpl().OfferAnswer(client, offerAnswer)
}

func ice_Call(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive ice_Call req")
	iceCall := &model.CSIceCall{}
	if err := mapstructure.Decode(data, iceCall); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	return impl.GetSigServiceImpl().IceCall(client, iceCall)
}

func ice_Answer(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive ice_Answer req")
	iceAnswer := &model.CSIceAnswer{}
	if err := mapstructure.Decode(data, iceAnswer); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	return impl.GetSigServiceImpl().IceAnswer(client, iceAnswer)
}

func media_ToAudio(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive media_ToAudio req")
	mediaToAudio := &model.CSMediaToAudio{}
	if err := mapstructure.Decode(data, mediaToAudio); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	return impl.GetSigServiceImpl().MediaToAudio(client, mediaToAudio)
}

func call_Mul(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive call_Mul req")
	callMul := &model.CSCallMul{}
	if err := mapstructure.Decode(data, callMul); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	callMulRsp := &model.ResponseStruct{}
	err := impl.GetSigServiceImpl().CallMul(client, callMul, callMulRsp)
	if err != nil {
		return constants.ErrCallMul
	}
	sendata, err := json.Marshal(callMulRsp)
	if err != nil {
		return constants.ErrJsonMarshal
	}

	err = client.SendMsg(sendata)
	if err != nil {
		return constants.ErrSendData
	}

	return nil

}

func call_MulAnswer(client *model.WSClient, data interface{}) error {

	logging.Logger.Info("receive call_MulAnswer req")

	callMulAnswer := &model.CSCallMulAnswer{}
	if err := mapstructure.Decode(data, callMulAnswer); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	callMulRsp := &model.ResponseStruct{}
	err := impl.GetSigServiceImpl().CallMulAnswer(client, callMulAnswer, callMulRsp)
	if err != nil {
		return constants.ErrCallMul
	}
	sendata, err := json.Marshal(callMulRsp)
	if err != nil {
		return constants.ErrJsonMarshal
	}

	err = client.SendMsg(sendata)
	if err != nil {
		return constants.ErrSendData
	}

	return nil
}

func call_MulRefuse(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive call_MulRefuse req")

	callMulRefuse := &model.CSCallMulRefuse{}
	if err := mapstructure.Decode(data, callMulRefuse); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	err := impl.GetSigServiceImpl().CallMulRefuse(client, callMulRefuse)
	if err != nil {
		return constants.ErrCallMul
	}

	return nil
}

func call_Hangup(client *model.WSClient, data interface{}) error {
	logging.Logger.Info("receive call_Hangup req")

	hangup := &model.CSMulHangup{}
	if err := mapstructure.Decode(data, hangup); err != nil {
		logging.Logger.Info(" map to struct failed, err is ", err)
		return constants.ErrMap2Struct
	}

	err := impl.GetSigServiceImpl().MulHangup(client, hangup)
	if err != nil {
		return constants.ErrCallMul
	}

	return nil
}
