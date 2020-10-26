package impl

import (
	"community-sig/config"
	"community-sig/constants"
	"community-sig/grpc_client"
	"community-sig/logging"
	"community-sig/model"
	"community-sig/protobuffer_def"
	"community-sig/service"
	"encoding/json"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type SigServiceImpl struct {
}

func (si *SigServiceImpl) UpdateHeartBeat(client *model.WSClient, heartReq *model.CSHeartBeat) error {

	ud := model.GetUserMgr().GetUser(heartReq.UserId)
	if ud == nil {
		logging.Logger.Infof("user %s not found ", heartReq.UserId)
		return constants.ErrUserNotFound
	}

	if ud.Client != client {
		logging.Logger.Info("handler client isn't equal to ud.Client ")
		return constants.ErrClient
	}

	ud.Client.UpdateHeartBeat()
	/*
		baseReq := &protobuffer_def.BaseRequest{}
		baseReq.RequestId = "update_on_line"
		baseReq.C = protobuffer_def.CMD_UPDATE_ON_LINE
		queryReq := &protobuffer_def.UpdateOnline{}
		queryReq.RegAddr = config.GetConf().Base.GRPCAddr
		queryReq.UserId = heartReq.UserId

		body, err := ptypes.MarshalAny(queryReq)
		if err != nil {
			logging.Logger.Info("proto marshal failed")
			return err
		}
		baseReq.Body = body

		logging.Logger.Info("baseReq is ", baseReq)
		grpc_client.SendRpcMsg(baseReq)
	*/
	return nil
}

func (si *SigServiceImpl) UserLogin(client *model.WSClient, loginReq *model.CSLogin,
	rsp *model.ResponseStruct) error {

	ud := &model.UserData{}
	ud.UserId = loginReq.UserId
	ud.Avator = loginReq.Avator
	ud.Phone = loginReq.Phone
	ud.RoomList = loginReq.RoomList
	ud.UserName = loginReq.UserName
	ud.Client = client
	model.GetUserMgr().AddUser(ud)

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "user_info_reg"
	baseReq.C = protobuffer_def.CMD_REGISTER_STATUS
	queryReq := &protobuffer_def.RegisterStatusRequest{}
	queryReq.Identity = loginReq.UserId
	queryReq.Phone = loginReq.Phone
	queryReq.RegisterInfo = config.GetConf().Base.GRPCAddr
	queryReq.RoomList = loginReq.RoomList
	queryReq.UserAvator = loginReq.Avator
	queryReq.UserName = loginReq.UserName

	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return err
	}
	baseReq.Body = body

	logging.Logger.Info("baseReq is ", baseReq)
	baseRsp, err := grpc_client.SendRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc request failed, err is ", err)
		rsp.Code = constants.ResponseCodeLoginFailed
		rsp.Desc = constants.ResponseCodeLoginFailed.String()
		rsp.Event = model.WS_LOGIN_SC
		rsp.Data = &model.SCLogin{UserId: loginReq.UserId}
		return err
	}

	if baseRsp == nil {
		rsp.Code = constants.ResponseCodeLoginFailed
		rsp.Desc = constants.ResponseCodeLoginFailed.String()
		rsp.Event = model.WS_LOGIN_SC
		rsp.Data = &model.SCLogin{UserId: loginReq.UserId}
		return nil
	}
	rsp.Code = constants.ResponseCodeSuccess
	rsp.Desc = "success"
	rsp.Event = model.WS_LOGIN_SC
	rsp.Data = &model.SCLogin{UserId: loginReq.UserId}
	return nil
}

//从status服务器获取用户信息
func (si *SigServiceImpl) RPCGetUserData(userId string) (*protobuffer_def.QueryStatusResponse, error) {
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "query_user_info"
	baseReq.C = protobuffer_def.CMD_QUERY_STATUS
	queryReq := &protobuffer_def.QueryStatusRequest{}
	queryReq.Identity = userId

	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil, err
	}
	baseReq.Body = body

	logging.Logger.Info("baseReq is ", baseReq)
	baseRsp, err := grpc_client.SendRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc request failed, err is ", err)
		return nil, err
	}

	if baseRsp.GetBody() == nil {
		logging.Logger.Info("user not found by rpc, err is ", err)
		return nil, constants.ErrUserNotFound
	}

	queryRsp := &protobuffer_def.QueryStatusResponse{}
	err = ptypes.UnmarshalAny(baseRsp.GetBody(), queryRsp)
	if err != nil {
		logging.Logger.Info("rpc request failed, err is ", err)
		return nil, err
	}

	logging.Logger.Info("get user data is ", queryRsp)

	return queryRsp, nil
}

func ChatForbid(lastChat int64) bool {
	cur := time.Now().Unix()
	if (cur - lastChat) < 3600 {
		return true
	}
	return false
}

//调用rpc，在status服务器上创建房间

func (si *SigServiceImpl) RPCCreateChatRoom(callerId string, answerId string) (*protobuffer_def.CreateChatRoomRsp, error) {
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "create_chat_room"
	baseReq.C = protobuffer_def.CMD_CREATE_CHAT_ROOM
	queryReq := &protobuffer_def.CreateChatRoomReq{}
	queryReq.Caller = callerId
	queryReq.Answer = answerId

	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil, err
	}
	baseReq.Body = body

	logging.Logger.Info("baseReq is ", baseReq)
	baseRsp, err := grpc_client.SendRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc request failed, err is ", err)
		return nil, err
	}

	if baseRsp.GetBody() == nil {
		logging.Logger.Info("create chat room failed ", err)
		return nil, constants.ErrRpcCreateChatRoom
	}

	queryRsp := &protobuffer_def.CreateChatRoomRsp{}
	err = ptypes.UnmarshalAny(baseRsp.GetBody(), queryRsp)
	if err != nil {
		logging.Logger.Info("rpc request failed, err is ", err)
		return nil, err
	}

	logging.Logger.Info("create chat room  data is ", queryRsp)

	return queryRsp, nil
}

func (si *SigServiceImpl) CallSingle(client *model.WSClient, callReq *model.CSCallSingle,
	rsp *model.ResponseStruct) error {

	//rpc获取
	callerrsp, err := si.RPCGetUserData(callReq.CallerId)
	if err != nil {
		rsp.Code = constants.ResponseCodeRpcGetUserFailed
		rsp.Desc = constants.ResponseCodeRpcGetUserFailed.String()
		rsp.Event = model.WS_CALL_SINGLE_SC
		return nil
	}

	answerrsp, err := si.RPCGetUserData(callReq.AnswerId)
	if err != nil {
		rsp.Code = constants.ResponseCodeRpcGetUserFailed
		rsp.Desc = constants.ResponseCodeRpcGetUserFailed.String()
		rsp.Event = model.WS_CALL_SINGLE_SC
		return nil
	}

	//判断对方是否在线
	if answerrsp.OffLine {
		rsp.Code = constants.ResponseCodeOnlineError
		rsp.Desc = "user isn't online  "
		rsp.Event = model.WS_CALL_SINGLE_SC
		return nil
	}

	if answerrsp.GetState() == constants.User_Busy && ChatForbid(answerrsp.LastChat) {
		rsp.Code = constants.ResponseUserBusyError
		rsp.Desc = rsp.Code.String()
		rsp.Event = model.WS_CALL_SINGLE_SC
		return nil
	}

	//将两个人关联到房间里，并设置busy信息
	creatrsp, err := si.RPCCreateChatRoom(callerrsp.Identity, callReq.AnswerId)
	if err != nil {
		rsp.Code = constants.ResponseUserBusyError
		rsp.Desc = rsp.Code.String()
		rsp.Event = model.WS_CALL_SINGLE_SC
		return nil
	}

	rsp.Code = constants.ResponseCodeSuccess
	rsp.Desc = "success"
	rsp.Event = model.WS_CALL_SINGLE_SC
	rsp.Data = &model.SCCallSingle{ChatRoomId: creatrsp.ChatRoomId}

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(answerrsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "single_call_notify"
	baseReq.C = protobuffer_def.CMD_SINGLE_CALL_NOTIFY
	notifyCall := &protobuffer_def.SingleCallNotify{}
	notifyCall.AnswerId = callReq.AnswerId
	notifyCall.CallerAvator = callerrsp.UserAvator
	notifyCall.CallerId = callerrsp.Identity
	notifyCall.CallerName = callerrsp.UserName
	notifyCall.CallerPhone = callerrsp.Phone
	notifyCall.ChatRoomId = creatrsp.ChatRoomId
	notifyCall.DeviceModel = callReq.DeviceModel
	notifyCall.MediaType = callReq.MediaType

	body, err := ptypes.MarshalAny(notifyCall)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil
	}
	baseReq.Body = body

	logging.Logger.Info("single call notify req body is ", notifyCall)

	rpcRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("sigRPCClient.PostRpcMsg failed")
		return nil
	}

	if rpcRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc send single call notify failed, code is ", rpcRsp.Code)
		return nil
	}

	logging.Logger.Info("send single call notify success")
	return nil
}

//向被叫方发送呼叫通知
func (si *SigServiceImpl) CallSingleNotify(client *model.WSClient, req *protobuffer_def.SingleCallNotify,
	roomId string) error {

	notifyRsp := &model.ResponseStruct{}
	notifyRsp.Code = constants.ResponseCodeSuccess
	notifyRsp.Desc = "success"
	notifyRsp.Event = model.WS_CALL_SINGLE_NOTIFY

	notifyRsp.Data = &model.SCCallSingleNotify{CallerId: req.CallerId, MediaType: req.MediaType,
		CallerAvator: req.CallerAvator,
		CallerName:   req.CallerName,
		CallerPhone:  req.CallerPhone,
		ChatRoomId:   roomId,
		DeviceModel:  req.DeviceModel}

	answerdata, err := json.Marshal(notifyRsp)
	if err != nil {
		logging.Logger.Info("json marshal failed, err is  ", err)
		return constants.ErrJsonMarshal
	}

	err = client.SendMsg(answerdata)
	if err != nil {
		logging.Logger.Info("answer client send failed, err is ", err)
		return nil
	}

	return nil
}

func (si *SigServiceImpl) RPCGetChatRoomData(chatRoomId string) (*protobuffer_def.GetChatRoomRsp, error) {
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "get_chat_room"
	baseReq.C = protobuffer_def.CMD_GET_CHAT_ROOM_REQ
	queryReq := &protobuffer_def.GetChatRoomReq{}
	queryReq.ChatRoomId = chatRoomId
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil, err
	}
	baseReq.Body = body

	logging.Logger.Info("baseReq is ", baseReq)
	baseRsp, err := grpc_client.SendRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc get chat room failed, err is ", err)
		return nil, err
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc get chat room failed, rsp  is nil")
		return nil, constants.ErrRpcGetChatRoom
	}

	queryRsp := &protobuffer_def.GetChatRoomRsp{}
	err = ptypes.UnmarshalAny(baseRsp.GetBody(), queryRsp)
	if err != nil {
		logging.Logger.Info("rpc get chat room  failed, err is ", err)
		return nil, err
	}

	logging.Logger.Info("get chat room data is ", queryRsp)
	return queryRsp, nil
}

func (si *SigServiceImpl) RPCKickUser(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.KickPerson{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive kick person is ", request)

	notifyer := model.GetUserMgr().GetUser(request.UserId)
	if notifyer == nil {
		logging.Logger.Infof("BeConvertId  %s is not exists", request.UserId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("BeConvertId %s is not online", request.UserId)
		return nil
	}

	model.GetUserMgr().KickUser(request.UserId)

	return nil
}

//被叫方同意接听
func (si *SigServiceImpl) AnswerSingle(client *model.WSClient, req *model.CSAnswerSingle,
	rsp *model.ResponseStruct) error {

	chatRoomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}
	//从status服务器获取主叫方信息
	callerRsp, err := si.RPCGetUserData(chatRoomRsp.Caller)
	if err != nil {
		return nil
	}

	if callerRsp.OffLine == true {
		logging.Logger.Infof("user %v not online ", callerRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "single_answer_notify"
	baseReq.C = protobuffer_def.CMD_SINGLE_ANSWER_SIG_TO_SIG
	queryReq := &protobuffer_def.SingleAnswerSigToSig{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.CallerId = callerRsp.Identity
	queryReq.DeviceModel = req.DeviceModel
	queryReq.AnswerId = chatRoomRsp.Answer
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto single_answer_notify marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(callerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc single answer notify failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc single answer notify failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc send single answer notify failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("send single answer notify success")

	return nil
}

func (si *SigServiceImpl) RefuseSingle(client *model.WSClient, req *model.CSRefuseSingle,
	rsp *model.ResponseStruct) error {

	roomData, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//解散房间
	si.RPCDelChatRoom(roomData.ChatRoomId)
	callerRsp, err := si.RPCGetUserData(roomData.Caller)
	if err != nil {
		return nil
	}

	if callerRsp.OffLine == true {
		logging.Logger.Infof("caller %v not online ", callerRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "single_refuse_notify"
	baseReq.C = protobuffer_def.CMD_SINGLE_REFUSE_SIG_TO_SIG
	queryReq := &protobuffer_def.SingleRefuseSigToSig{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.CallerId = callerRsp.Identity
	queryReq.AnswerId = roomData.Answer
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto single_refuse_notify marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(callerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc single answer notify failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc single answer notify failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc send single answer notify failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("send single refuse notify success")

	return nil
}

func (si *SigServiceImpl) HangupSingle(client *model.WSClient, req *model.CSHangupSingle,
	rsp *model.ResponseStruct) error {

	ud := model.GetC2UMgr().GetUserByClient(client.Id())
	if ud == nil {
		logging.Logger.Infof("get user by client %s not exists", client.Id())
		return nil
	}

	roomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//校验是不是主叫方挂断
	if ud.UserId != roomRsp.Caller {
		logging.Logger.Infof("hang up is not caller, user is %s, caller is %s",
			ud.UserId, roomRsp.Caller)
		return constants.ErrHangupIsNotCaller
	}

	answerRsp, err := si.RPCGetUserData(roomRsp.Answer)
	if err != nil {
		return nil
	}

	if answerRsp.OffLine {
		logging.Logger.Infof("answer %s is not online", answerRsp.Identity)
		//删除聊天室
		si.RPCDelChatRoom(req.ChatRoomId)
		return nil
	}
	//删除聊天室
	si.RPCDelChatRoom(req.ChatRoomId)

	sigRPCClient, err := grpc_client.GetSigRPCClient(answerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "hang_up_notify"
	baseReq.C = protobuffer_def.CMD_HANG_UP_SIG_TO_SIG
	notifyCall := &protobuffer_def.HangUpSigToSig{}
	notifyCall.AnswerId = answerRsp.Identity
	notifyCall.ChatRoomId = req.ChatRoomId
	body, err := ptypes.MarshalAny(notifyCall)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil
	}
	baseReq.Body = body

	logging.Logger.Info("Hang up notify req body is ", notifyCall)

	rpcRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("sigRPCClient.PostRpcMsg failed")
		return nil
	}

	if rpcRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc send hang up notify failed, code is ", rpcRsp.Code)
		return nil
	}

	logging.Logger.Info("send hang up notify success")
	return nil
}

func (si *SigServiceImpl) RPCHangupNotify(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.HangUpSigToSig{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc hang up notify is ", request)

	answer := model.GetUserMgr().GetUser(request.AnswerId)
	if answer == nil {
		logging.Logger.Infof("answer %s is not exists", request.AnswerId)
		return nil
	}

	if !answer.IsOnline() {
		logging.Logger.Infof("answer %s is not online", request.AnswerId)
		return nil
	}

	hangupNotify := &model.SCHangupSingleNotify{ChatRoomId: request.ChatRoomId}
	hanguprsp := model.ResponseStruct{}
	hanguprsp.Code = constants.ResponseCodeSuccess
	hanguprsp.Desc = constants.ResponseCodeSuccess.String()
	hanguprsp.Data = hangupNotify
	hanguprsp.Event = model.WS_SINGLE_HANGUP_NOTIFY

	sendata, _ := json.Marshal(hanguprsp)
	answer.Client.SendMsg(sendata)

	return nil

}

func (si *SigServiceImpl) TerminalSingle(client *model.WSClient, req *model.CSTerminateSingle) error {

	ud := model.GetC2UMgr().GetUserByClient(client.Id())
	if ud == nil {
		logging.Logger.Infof("get user by client %s not exists", client.Id())
		return nil
	}

	roomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//校验挂断方是否和请求id相同

	if ud.UserId != req.CancelId {
		logging.Logger.Infof("terminate user is invalid, client user is %s ,req cancel is %s ",
			ud.UserId, req.CancelId)
		return constants.ErrTerminateUserInvalid
	}

	notifyId := roomRsp.Answer
	if req.CancelId == roomRsp.Answer {
		notifyId = roomRsp.Caller
	}

	notifyRsp, err := si.RPCGetUserData(notifyId)
	if err != nil {
		return nil
	}

	if notifyRsp.OffLine {
		logging.Logger.Infof("answer %s is not online", notifyRsp.Identity)
		//删除聊天室
		si.RPCDelChatRoom(req.ChatRoomId)
		return nil
	}
	//删除聊天室
	si.RPCDelChatRoom(req.ChatRoomId)

	sigRPCClient, err := grpc_client.GetSigRPCClient(notifyRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "single_terminal_notify"
	baseReq.C = protobuffer_def.CMD_TERMINAL_SIG_TO_SIG
	notifyCall := &protobuffer_def.SingleTerminal{}
	notifyCall.CancelId = req.CancelId
	notifyCall.ChatRoomId = req.ChatRoomId
	notifyCall.BeCanceledId = notifyId
	body, err := ptypes.MarshalAny(notifyCall)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil
	}
	baseReq.Body = body

	logging.Logger.Info("terminate call notify req body is ", notifyCall)

	rpcRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("sigRPCClient.PostRpcMsg failed")
		return nil
	}

	if rpcRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc send hang up notify failed, code is ", rpcRsp.Code)
		return nil
	}

	logging.Logger.Info("terminate call notify success")
	return nil

}

func (si *SigServiceImpl) RPCTerminateNotify(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.SingleTerminal{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc single terminate notify is ", request)

	notifyer := model.GetUserMgr().GetUser(request.BeCanceledId)
	if notifyer == nil {
		logging.Logger.Infof("becanceled  %s is not exists", request.BeCanceledId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("becanceled %s is not online", request.BeCanceledId)
		return nil
	}

	terminalNotify := &model.SCTerminateSingleNotify{ChatRoomId: request.ChatRoomId,
		CancelId: request.CancelId}
	terminalNotifyRsp := model.ResponseStruct{}
	terminalNotifyRsp.Code = constants.ResponseCodeSuccess
	terminalNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	terminalNotifyRsp.Data = terminalNotify
	terminalNotifyRsp.Event = model.WS_SINGLE_TERMINATE_NoTIFY

	sendata, _ := json.Marshal(terminalNotifyRsp)
	notifyer.Client.SendMsg(sendata)

	return nil

}

func (si *SigServiceImpl) OfferCall(client *model.WSClient, req *model.CSOfferCall) error {

	chatRoomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//判断参数是否合理
	if req.CallerId != chatRoomRsp.Caller {
		logging.Logger.Infof("caller %s is invalid, chat room caller is %s", req.CallerId,
			chatRoomRsp.Caller)
		return nil
	}

	//从status服务器获取被叫方信息
	answerRsp, err := si.RPCGetUserData(chatRoomRsp.Answer)
	if err != nil {
		return nil
	}

	if answerRsp.OffLine == true {
		logging.Logger.Infof("user %v not online ", answerRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "offer_call"
	baseReq.C = protobuffer_def.CMD_OFFER_CALL_SIG_TO_SIG
	queryReq := &protobuffer_def.OfferCall{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.CallerId = req.CallerId
	queryReq.AnswerId = chatRoomRsp.Answer
	queryReq.Sdp = req.Sdp
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto call offer marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(answerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc offer call failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc offer call  failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc offer call  failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("send offer call success")

	return nil

}

func (si *SigServiceImpl) RPCOfferCall(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.OfferCall{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc offer call is ", request)

	notifyer := model.GetUserMgr().GetUser(request.AnswerId)
	if notifyer == nil {
		logging.Logger.Infof("answer  %s is not exists", request.AnswerId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("answer %s is not online", request.AnswerId)
		return nil
	}

	offCallNotify := &model.SCOfferCallNotify{ChatRoomId: request.ChatRoomId,
		CallerId: request.CallerId, Sdp: request.Sdp}
	offCallNotifyRsp := model.ResponseStruct{}
	offCallNotifyRsp.Code = constants.ResponseCodeSuccess
	offCallNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	offCallNotifyRsp.Data = offCallNotify
	offCallNotifyRsp.Event = model.WS_OFFER_CALL_NOTIFY

	sendata, _ := json.Marshal(offCallNotifyRsp)
	notifyer.Client.SendMsg(sendata)

	return nil
}

func (si *SigServiceImpl) OfferAnswer(client *model.WSClient, req *model.CSOfferAnswer) error {

	chatRoomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//判断参数是否合理
	if req.AnswerId != chatRoomRsp.Answer {
		logging.Logger.Infof("answer %s is invalid, chat room answer is %s", req.AnswerId,
			chatRoomRsp.Answer)
		return nil
	}

	//从status服务器获取主叫方信息
	callerRsp, err := si.RPCGetUserData(chatRoomRsp.Caller)
	if err != nil {
		return nil
	}

	if callerRsp.OffLine == true {
		logging.Logger.Infof("user %v not online ", callerRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "offer_answer"
	baseReq.C = protobuffer_def.CMD_OFFER_ANSWER_SIG_TO_SIG
	queryReq := &protobuffer_def.OfferAnswer{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.CallerId = chatRoomRsp.Caller
	queryReq.AnswerId = chatRoomRsp.Answer
	queryReq.Sdp = req.Sdp
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto call offer marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(callerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc offer answer failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc offer answer  failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc offer answer  failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("send offer answer success")

	return nil

}

//服务器收到answer端服务器发送的offer
func (si *SigServiceImpl) RPCOfferAnswer(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.OfferAnswer{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc offer answer is ", request)

	notifyer := model.GetUserMgr().GetUser(request.CallerId)
	if notifyer == nil {
		logging.Logger.Infof("caller  %s is not exists", request.CallerId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("caller %s is not online", request.CallerId)
		return nil
	}

	offAnswerNotify := &model.SCOfferAnswerNotify{ChatRoomId: request.ChatRoomId,
		AnswerId: request.AnswerId, Sdp: request.Sdp}
	offAnswerNotifyRsp := model.ResponseStruct{}
	offAnswerNotifyRsp.Code = constants.ResponseCodeSuccess
	offAnswerNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	offAnswerNotifyRsp.Data = offAnswerNotify
	offAnswerNotifyRsp.Event = model.WS_OFFER_ANSWER_NOTIFY

	sendata, _ := json.Marshal(offAnswerNotifyRsp)
	notifyer.Client.SendMsg(sendata)

	return nil
}

func (si *SigServiceImpl) IceCall(client *model.WSClient, req *model.CSIceCall) error {
	chatRoomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//判断参数是否合理
	if req.CallerId != chatRoomRsp.Caller {
		logging.Logger.Infof("caller %s is invalid, chat room caller is %s", req.CallerId,
			chatRoomRsp.Caller)
		return nil
	}

	//从status服务器获取被叫方信息
	answerRsp, err := si.RPCGetUserData(chatRoomRsp.Answer)
	if err != nil {
		return nil
	}

	if answerRsp.OffLine == true {
		logging.Logger.Infof("user %v not online ", answerRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "ice_call"
	baseReq.C = protobuffer_def.CMD_ICE_CALL_SIG_TO_SIG
	queryReq := &protobuffer_def.IceCall{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.CallerId = req.CallerId
	queryReq.AnswerId = chatRoomRsp.Answer
	queryReq.IceCandidate = req.IceCandidate
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto call ice marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(answerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc ice call failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc ice call  failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc ice call  failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("send ice call success")

	return nil

}

func (si *SigServiceImpl) RPCIceCall(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.IceCall{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc ice call is ", request)

	notifyer := model.GetUserMgr().GetUser(request.AnswerId)
	if notifyer == nil {
		logging.Logger.Infof("answer  %s is not exists", request.AnswerId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("answer %s is not online", request.AnswerId)
		return nil
	}

	iceCallNotify := &model.SCIceCallNotify{ChatRoomId: request.ChatRoomId, CallerId: request.CallerId,
		IceCandidate: request.IceCandidate}
	iceCallNotifyRsp := model.ResponseStruct{}
	iceCallNotifyRsp.Code = constants.ResponseCodeSuccess
	iceCallNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	iceCallNotifyRsp.Data = iceCallNotify
	iceCallNotifyRsp.Event = model.WS_ICE_CALL_NOTIFY

	sendata, _ := json.Marshal(iceCallNotifyRsp)
	notifyer.Client.SendMsg(sendata)
	return nil
}

func (si *SigServiceImpl) IceAnswer(client *model.WSClient, req *model.CSIceAnswer) error {
	chatRoomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	//判断参数是否合理
	if req.AnswerId != chatRoomRsp.Answer {
		logging.Logger.Infof("answer %s is invalid, chat room answer is %s", req.AnswerId,
			chatRoomRsp.Answer)
		return nil
	}

	//从status服务器获取主叫方信息
	callerRsp, err := si.RPCGetUserData(chatRoomRsp.Caller)
	if err != nil {
		return nil
	}

	if callerRsp.OffLine == true {
		logging.Logger.Infof("user %v not online ", callerRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "ice_answer"
	baseReq.C = protobuffer_def.CMD_ICE_ANSWER_SIG_TO_SIG
	queryReq := &protobuffer_def.IceAnswer{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.CallerId = chatRoomRsp.Caller
	queryReq.AnswerId = chatRoomRsp.Answer
	queryReq.IceCandidate = req.IceCandidate
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto answer ice marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(callerRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc ice answer failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc ice answer  failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("rpc ice answer  failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("send ice answer success")

	return nil

}

func (si *SigServiceImpl) RPCIceAnswer(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.IceAnswer{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc ice answer is ", request)

	notifyer := model.GetUserMgr().GetUser(request.CallerId)
	if notifyer == nil {
		logging.Logger.Infof("caller  %s is not exists", request.CallerId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("caller %s is not online", request.CallerId)
		return nil
	}

	iceAnswerNotify := &model.SCIceAnswerNotify{ChatRoomId: request.ChatRoomId, AnswerId: request.AnswerId,
		IceCandidate: request.IceCandidate}
	iceAnswerNotifyRsp := model.ResponseStruct{}
	iceAnswerNotifyRsp.Code = constants.ResponseCodeSuccess
	iceAnswerNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	iceAnswerNotifyRsp.Data = iceAnswerNotify
	iceAnswerNotifyRsp.Event = model.WS_ICE_ANSWER_NOTIFY

	sendata, _ := json.Marshal(iceAnswerNotifyRsp)
	notifyer.Client.SendMsg(sendata)

	return nil
}

func (si *SigServiceImpl) MediaToAudio(client *model.WSClient, req *model.CSMediaToAudio) error {

	chatRoomRsp, err := si.RPCGetChatRoomData(req.ChatRoomId)
	if err != nil {
		return nil
	}

	ud := model.GetC2UMgr().GetUserByClient(client.Id())
	if ud == nil {
		logging.Logger.Infof("get user by client %s not exists", client.Id())
		return nil
	}

	//校验是不是转换人
	if ud.UserId != req.ConverId {
		logging.Logger.Infof("media to audio is invalid, client user is %s ,req conver is %s ",
			ud.UserId, req.ConverId)
		return constants.ErrMediaToAudioInvalid
	}

	notifyId := chatRoomRsp.Answer
	if req.ConverId == chatRoomRsp.Answer {
		notifyId = chatRoomRsp.Caller
	}

	//从status服务器获取被通知方信息
	notifyRsp, err := si.RPCGetUserData(notifyId)
	if err != nil {
		return nil
	}

	if notifyRsp.OffLine == true {
		logging.Logger.Infof("user %v not online ", notifyRsp.Identity)
		return nil
	}

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "media_to_audio"
	baseReq.C = protobuffer_def.CMD_CMD_MEDIA_TO_AUDIO
	queryReq := &protobuffer_def.MediaToAudio{}
	queryReq.ChatRoomId = req.ChatRoomId
	queryReq.ConvertId = req.ConverId
	queryReq.BeConvertId = notifyRsp.Identity
	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("media to audio marshal failed")
		return nil
	}
	baseReq.Body = body

	//调用rpc通知另一方呼叫消息，因为另一方可能不在同服务器中
	//从连接池获取对端连接的客户端

	sigRPCClient, err := grpc_client.GetSigRPCClient(notifyRsp.RegAddr)
	if err != nil {
		logging.Logger.Info("get sig rpc client failed, err is ", err)
		return nil
	}

	baseRsp, err := sigRPCClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("media to audio failed, err is ", err)
		return nil
	}

	if baseRsp == nil {
		logging.Logger.Info("media to audio  failed, rsp  is nil")
		return nil
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("media to audio  failed, code is ", baseRsp.Code)
		return nil
	}

	logging.Logger.Info("media to audio success")

	return nil

}

func (si *SigServiceImpl) RPCMediaToAudio(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.MediaToAudio{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc media to audio is ", request)

	notifyer := model.GetUserMgr().GetUser(request.BeConvertId)
	if notifyer == nil {
		logging.Logger.Infof("BeConvertId  %s is not exists", request.BeConvertId)
		return nil
	}

	if !notifyer.IsOnline() {
		logging.Logger.Infof("BeConvertId %s is not online", request.BeConvertId)
		return nil
	}

	mediaToAudioNotify := &model.SCMediaToAudioNotify{ChatRoomId: request.ChatRoomId,
		ConverId: request.ConvertId}
	mediaToAudioNotifyRsp := model.ResponseStruct{}
	mediaToAudioNotifyRsp.Code = constants.ResponseCodeSuccess
	mediaToAudioNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	mediaToAudioNotifyRsp.Data = mediaToAudioNotify
	mediaToAudioNotifyRsp.Event = model.WS_MEDIA_TO_AUDIO_NOTIFY

	sendata, _ := json.Marshal(mediaToAudioNotifyRsp)
	notifyer.Client.SendMsg(sendata)

	return nil
}

func (si *SigServiceImpl) CallMul(client *model.WSClient, callReq *model.CSCallMul, callRsp *model.ResponseStruct) error {

	rd := model.GetRoomMgr().GetRoom(callReq.RoomId)
	if rd == nil {
		callRsp.Code = constants.ResponseCodeRoomError
		callRsp.Desc = constants.ResponseCodeRoomError.String()
		callRsp.Event = model.WS_CALL_MULT_SC
		logging.Logger.Infof("room id %s not found", callReq.RoomId)
		return nil
	}

	ud := model.GetC2UMgr().GetUserByClient(client.Id())
	if ud == nil {
		callRsp.Code = constants.ResponseCodeRoomError
		callRsp.Desc = constants.ResponseCodeRoomError.String()
		callRsp.Event = model.WS_CALL_MULT_SC
		logging.Logger.Infof("get user by client %s not exists", client.Id())
		return nil
	}

	if ud.UserId != callReq.CallerId {
		callRsp.Code = constants.ResponseCodeRoomError
		callRsp.Desc = constants.ResponseCodeRoomError.String()
		callRsp.Event = model.WS_CALL_MULT_SC
		logging.Logger.Infof("caller %s invalid ", ud.UserId)
		return nil
	}

	if ud.GetState() == constants.User_Busy && ud.ChatForbid() {
		callRsp.Code = constants.ResponseCodeOnlineError
		callRsp.Desc = constants.ResponseCodeOnlineError.String()
		callRsp.Event = model.WS_CALL_MULT_SC
		logging.Logger.Infof("caller %s is forbiden call ", ud.UserId)
		return nil
	}

	//判断房间内人是否在线
	allOffLine := true
	onlineUserMap := make(map[string]*model.UserData)
	for _, roomUd := range rd.UserMap {
		if !roomUd.IsOnline() {
			continue
		}

		if roomUd.GetState() == constants.User_Busy && roomUd.ChatForbid() {
			continue
		}

		onlineUserMap[roomUd.UserId] = roomUd
		allOffLine = false
	}

	if allOffLine == true {
		callRsp.Code = constants.ResponseRoomUserALLOff
		callRsp.Desc = constants.ResponseRoomUserALLOff.String()
		callRsp.Event = model.WS_CALL_MULT_SC
		logging.Logger.Infof("room %s all user offline or busy ", callReq.RoomId)
		return nil
	}

	mulChatRoom := model.GetMulChatMgr().CreateMulChat(callReq.CallerId, callReq.RoomId, onlineUserMap)

	callMulNotify := &model.SCCallMulNotify{CallerId: callReq.CallerId, ChatRoomId: mulChatRoom.ID(),
		MediaType: callReq.RoomId, DeviceModel: callReq.DeviceModel}
	callMulNotifyRsp := model.ResponseStruct{}
	callMulNotifyRsp.Code = constants.ResponseCodeSuccess
	callMulNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	callMulNotifyRsp.Data = callMulNotify
	callMulNotifyRsp.Event = model.WS_CALL_MULT_NOTIFY

	sendata, _ := json.Marshal(callMulNotifyRsp)

	//通知在线用户来电
	for _, onlineUd := range onlineUserMap {

		onlineUd.Client.SendMsg(sendata)
	}

	callRsp.Code = constants.ResponseCodeSuccess
	callRsp.Desc = constants.ResponseCodeSuccess.String()
	callRsp.Event = model.WS_CALL_MULT_SC
	mulCallSC := &model.SCCallMul{}
	mulCallSC.ChatRoomId = mulChatRoom.ID()
	callRsp.Data = mulCallSC

	return nil
}

func (si *SigServiceImpl) CallMulAnswer(client *model.WSClient, answerReq *model.CSCallMulAnswer,
	answerRsp *model.ResponseStruct) error {

	mulChatRoom := model.GetMulChatMgr().GetMulChatRoom(answerReq.ChatRoomId)
	if mulChatRoom == nil {
		answerRsp.Code = constants.ResponseMulChatRoomError
		answerRsp.Desc = constants.ResponseMulChatRoomError.String()
		answerRsp.Event = model.WS_CALL_MULT_ANSWER_SC
		logging.Logger.Infof("mul chat room %s not exist", answerReq.ChatRoomId)
		return nil
	}

	answerValid := mulChatRoom.AnswerValid(answerReq.AnswerId)
	if !answerValid {
		answerRsp.Code = constants.ResponseMulChatAnswerError
		answerRsp.Desc = constants.ResponseMulChatAnswerError.String()
		answerRsp.Event = model.WS_CALL_MULT_ANSWER_SC
		logging.Logger.Infof("answer %s is not in mul chat room ", answerReq.AnswerId)
		return nil
	}

	chatRoomId := mulChatRoom.ID()
	callerId := mulChatRoom.Caller()
	answerId := answerReq.AnswerId

	callerUd := model.GetUserMgr().GetUser(callerId)
	if !callerUd.IsOnline() {
		answerRsp.Code = constants.ResponseUserNotOnline
		answerRsp.Desc = constants.ResponseUserNotOnline.String()
		answerRsp.Event = model.WS_CALL_MULT_ANSWER_SC
		logging.Logger.Infof("caller %s is not online", callerId)
		return nil
	}

	model.GetMulChatMgr().DelMulRoomWaiters(chatRoomId, answerId)
	model.GetChatMgr().Mul2ChatRoom(chatRoomId, callerId, answerId)

	//向主叫方通知接听结果

	answerMulNotify := &model.SCCallMulAnswerNotify{ChatRoomId: answerReq.ChatRoomId,
		DeviceModel: answerReq.DeviceModel}
	answerMulNotifyRsp := model.ResponseStruct{}
	answerMulNotifyRsp.Code = constants.ResponseCodeSuccess
	answerMulNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	answerMulNotifyRsp.Data = answerMulNotify
	answerMulNotifyRsp.Event = model.WS_CALL_MULT_ANSWER_NOTIFY

	sendata, _ := json.Marshal(answerMulNotify)
	callerUd.Client.SendMsg(sendata)
	// 给被叫方回复信息

	answerRsp.Code = constants.ResponseCodeSuccess
	answerRsp.Desc = constants.ResponseCodeSuccess.String()
	answerRsp.Event = model.WS_CALL_MULT_ANSWER_SC
	mulCallSC := &model.SCCallMulAnswer{}
	mulCallSC.ChatRoomId = answerReq.ChatRoomId

	answerRsp.Data = mulCallSC

	return nil
}

func (si *SigServiceImpl) CallMulRefuse(client *model.WSClient, refuseReq *model.CSCallMulRefuse) error {
	refuseUd := model.GetC2UMgr().GetUserByClient(client.Id())
	if refuseUd == nil {
		return nil
	}
	mulChatRoom := model.GetMulChatMgr().GetMulChatRoom(refuseReq.ChatRoomId)
	if mulChatRoom == nil {
		return nil
	}

	callUd := model.GetUserMgr().GetUser(mulChatRoom.Caller())
	if callUd == nil {
		return nil
	}

	if mulChatRoom.DeleteWaitAnswer(refuseUd.UserId) {
		refuseMulNotify := &model.SCCallMulRefuseNotify{ChatRoomId: refuseReq.ChatRoomId}
		refuseMulNotifyRsp := model.ResponseStruct{}
		refuseMulNotifyRsp.Code = constants.ResponseCodeSuccess
		refuseMulNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
		refuseMulNotifyRsp.Data = refuseMulNotify
		refuseMulNotifyRsp.Event = model.WS_CALL_MULT_REFUSE_NOTIFY

		sendata, _ := json.Marshal(refuseMulNotifyRsp)
		callUd.Client.SendMsg(sendata)
	}

	return nil
}

func (si *SigServiceImpl) MulHangup(client *model.WSClient, refuseReq *model.CSMulHangup) error {

	callerUd := model.GetC2UMgr().GetUserByClient(client.Id())
	if callerUd == nil {
		return nil
	}
	mulChatRoom := model.GetMulChatMgr().GetMulChatRoom(refuseReq.ChatRoomId)
	if mulChatRoom == nil {
		return nil
	}

	callUd := model.GetUserMgr().GetUser(mulChatRoom.Caller())
	if callUd == nil {
		return nil
	}

	model.GetMulChatMgr().NotifyAnswerHangup(refuseReq.ChatRoomId)

	return nil
}

func (si *SigServiceImpl) preDeal(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse, request proto.Message) bool {
	//解析body
	if baseRequest.GetBody() == nil {
		baseResponse.Code = protobuffer_def.ReturnCode_BODY_IS_NULL
		baseResponse.Desc = "body is null"
		return false
	}
	//反序列化
	err := ptypes.UnmarshalAny(baseRequest.GetBody(), request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "body deserialization error"
		return false
	}
	return true
}

//rpc server 收到 主叫方服务器发送的callsinglenotify请求
func (si *SigServiceImpl) RPCCallSingleNotify(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS

	//解析请求参数
	request := &protobuffer_def.SingleCallNotify{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc call single notify is ", request)
	answer := model.GetUserMgr().GetUser(request.AnswerId)
	if answer == nil {
		logging.Logger.Info("get user failed, userid is ", request.AnswerId)
		return nil
	}

	if answer.IsOnline() == false {
		logging.Logger.Info("answer is not online, userid is ", request.AnswerId)
		return nil
	}
	si.CallSingleNotify(answer.Client, request, request.ChatRoomId)

	return nil
}

func (si *SigServiceImpl) RPCForceTerminateNotify(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.ForceTerminateNotify{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc force terminate notify is ", request)
	other := model.GetUserMgr().GetUser(request.OtherId)
	if other == nil {
		logging.Logger.Info("get user failed, userid is ", request.OtherId)
		return nil
	}

	if other.IsOnline() == false {
		logging.Logger.Info("other is not online, userid is ", request.OtherId)
		return nil
	}

	terminalNotify := &model.SCForceTerminateNotify{ChatRoomId: request.ChatRoomId}

	terminalRsp := model.ResponseStruct{}
	terminalRsp.Event = model.WS_FORCE_TERMINATE_NOTIFY
	terminalRsp.Code = constants.ResponseCodeSuccess
	terminalRsp.Data = terminalNotify
	terminalRsp.Desc = constants.ResponseCodeSuccess.String()
	sendata, err := json.Marshal(terminalRsp)
	if err != nil {
		logging.Logger.Info("terminate json marshal failed")
		return nil
	}
	other.Client.SendMsg(sendata)

	return nil
}

//主叫方服务器收到通知，将应答结果推送给客户端
func (si *SigServiceImpl) RPCSingleAnswerNotify(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.SingleAnswerSigToSig{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc single answer notify is ", request)

	answerRsp, err := si.RPCGetUserData(request.AnswerId)
	if err != nil {
		logging.Logger.Info("get answer from status failed")
		return nil
	}

	caller := model.GetUserMgr().GetUser(request.CallerId)
	if caller == nil {
		logging.Logger.Info("get user failed, userid is ", request.CallerId)
		return nil
	}

	if caller.IsOnline() == false {
		logging.Logger.Info("other is not online, userid is ", caller.UserId)
		return nil
	}

	//通知主叫方
	anserNotify := &model.SCAnswerSingleNotify{AnswerId: request.AnswerId, ChatRoomId: request.ChatRoomId,
		AnswerAvator: answerRsp.UserAvator,
		AnswerName:   answerRsp.UserName,
		AnswerPhone:  answerRsp.Phone,
		DeviceModel:  request.DeviceModel}
	answerNotifyRsp := model.ResponseStruct{}
	answerNotifyRsp.Event = model.WS_CALL_SINGLE_ANSWER_NOTIFY
	answerNotifyRsp.Code = constants.ResponseCodeSuccess
	answerNotifyRsp.Data = anserNotify
	answerNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	sendata, _ := json.Marshal(answerNotifyRsp)
	caller.Client.SendMsg(sendata)
	return nil
}

func (si *SigServiceImpl) RPCDelChatRoom(chatRoomId string) error {

	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "del_chat_room"
	baseReq.C = protobuffer_def.CMD_DEL_CHAT_ROOM
	queryReq := &protobuffer_def.DelChatRoomReq{}
	queryReq.ChatRoomId = chatRoomId

	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return err
	}
	baseReq.Body = body

	baseRsp, err := grpc_client.SendRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("rpc del chat room failed, err is ", err)
		return err
	}

	if baseRsp == nil {
		logging.Logger.Info("rpc del chat room failed, rsp is nil")
		return nil
	}
	return nil
}

//主叫方服务器收到被叫方服务器发送的拒绝通知
func (si *SigServiceImpl) RPCSingleRefuseNotify(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.SingleRefuseSigToSig{}
	if !si.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "deserialize failed"
		logging.Logger.Info("json unmarshal failed")
		return nil
	}

	logging.Logger.Info("receive rpc single refuse notify is ", request)

	caller := model.GetUserMgr().GetUser(request.CallerId)
	if caller == nil {
		logging.Logger.Info("get user failed, userid is ", request.CallerId)
		return nil
	}

	if caller.IsOnline() == false {
		logging.Logger.Info("other is not online, userid is ", caller.UserId)
		return nil
	}

	//通知主叫方
	refuseNotify := &model.SCRefuseSingleNotify{AnswerId: request.AnswerId, ChatRoomId: request.ChatRoomId}
	refuseNotifyRsp := model.ResponseStruct{}
	refuseNotifyRsp.Event = model.WS_CALL_SINGLE_REFUSE_NOTIFY
	refuseNotifyRsp.Code = constants.ResponseCodeSuccess
	refuseNotifyRsp.Data = refuseNotify
	refuseNotifyRsp.Desc = constants.ResponseCodeSuccess.String()
	sendata, _ := json.Marshal(refuseNotifyRsp)
	caller.Client.SendMsg(sendata)

	return nil

}

var sigServiceImpl service.SigService
var sigServiceOnce sync.Once

func GetSigServiceImpl() service.SigService {
	sigServiceOnce.Do(func() {
		sigServiceImpl = &SigServiceImpl{}
	})
	return sigServiceImpl
}
