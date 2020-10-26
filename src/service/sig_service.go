package service

import (
	"community-sig/model"
	"community-sig/protobuffer_def"
)

type SigService interface {
	UpdateHeartBeat(*model.WSClient, *model.CSHeartBeat) error
	UserLogin(*model.WSClient, *model.CSLogin, *model.ResponseStruct) error
	CallSingle(*model.WSClient, *model.CSCallSingle, *model.ResponseStruct) error
	CallSingleNotify(*model.WSClient, *protobuffer_def.SingleCallNotify, string) error
	AnswerSingle(*model.WSClient, *model.CSAnswerSingle, *model.ResponseStruct) error
	RefuseSingle(*model.WSClient, *model.CSRefuseSingle, *model.ResponseStruct) error
	HangupSingle(*model.WSClient, *model.CSHangupSingle, *model.ResponseStruct) error
	TerminalSingle(*model.WSClient, *model.CSTerminateSingle) error
	OfferCall(*model.WSClient, *model.CSOfferCall) error
	OfferAnswer(*model.WSClient, *model.CSOfferAnswer) error
	IceCall(*model.WSClient, *model.CSIceCall) error
	IceAnswer(*model.WSClient, *model.CSIceAnswer) error
	MediaToAudio(*model.WSClient, *model.CSMediaToAudio) error
	CallMul(*model.WSClient, *model.CSCallMul, *model.ResponseStruct) error
	CallMulAnswer(*model.WSClient, *model.CSCallMulAnswer, *model.ResponseStruct) error
	CallMulRefuse(*model.WSClient, *model.CSCallMulRefuse) error
	MulHangup(*model.WSClient, *model.CSMulHangup) error

	// rpc-接口
	RPCCallSingleNotify(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务接收强制中止会话逻辑
	RPCForceTerminateNotify(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器接收到被叫方服务器发送的接听通知
	RPCSingleAnswerNotify(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器接收到被叫方服务器发送的拒绝通知
	RPCSingleRefuseNotify(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到主叫方服务器挂断通知
	RPCHangupNotify(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到某一方服务器中断通话通知
	RPCTerminateNotify(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到caller服务器发送的offer通知
	RPCOfferCall(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到answer服务器发送的offer通知
	RPCOfferAnswer(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到caller服务器发送的ice通知
	RPCIceCall(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到answer服务器发送的ice通知
	RPCIceAnswer(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到另一方服务器发送的视频转语音通知
	RPCMediaToAudio(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	//rpc 服务器收到status服务器的下线请求
	RPCKickUser(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
}
