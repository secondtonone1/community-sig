package grpc_server

import (
	"community-sig/config"
	"community-sig/protobuffer_def"
	"community-sig/service/impl"
	"context"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"
)

var (
	statusServers service.Service
	statusOnce    = &sync.Once{}
	statusServer  server.Server
)

//初始化grpc service服务
func StartComSigService(config *config.Config) {
	statusOnce.Do(func() {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// create GRPC service
		service := grpc.NewService(
			service.Address(config.Base.GRPCAddr),
			service.Name(config.Base.ServiceName),
			service.Registry(config.RegisterCenter.GetRegisterCenter()),
			service.RegisterTTL(time.Second*30),
			service.RegisterInterval(time.Second*20),
			service.Context(ctx),
		)

		service.Init()

		statusServer = service.Server()
		// register test handler
		protobuffer_def.RegisterComSigServerHandler(service.Server(), &comSigGrpcServiceImpl{})

		//启动服务
		if err := service.Run(); err != nil {
			panic(err)
		}
	})
}

func StopComSigService() {
	statusServer.Stop()
}

type comSigGrpcServiceImpl struct{}

func (s *comSigGrpcServiceImpl) BaseInterface(context context.Context, baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.C = baseRequest.GetC()
	baseResponse.RequestId = baseRequest.GetRequestId()
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS

	switch baseRequest.GetC() {

	case protobuffer_def.CMD_SINGLE_CALL_NOTIFY: //收到另一端的呼叫通知
		return impl.GetSigServiceImpl().RPCCallSingleNotify(baseRequest, baseResponse)
	case protobuffer_def.CMD_FORCE_TERMINAL_NOTIFY: //服务器接收到强制终止通话逻辑
		return impl.GetSigServiceImpl().RPCForceTerminateNotify(baseRequest, baseResponse)
	case protobuffer_def.CMD_SINGLE_ANSWER_SIG_TO_SIG: //接收到answer服务器发送的接听通知
		return impl.GetSigServiceImpl().RPCSingleAnswerNotify(baseRequest, baseResponse)
	case protobuffer_def.CMD_SINGLE_REFUSE_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCSingleRefuseNotify(baseRequest, baseResponse)
	case protobuffer_def.CMD_HANG_UP_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCHangupNotify(baseRequest, baseResponse)
	case protobuffer_def.CMD_TERMINAL_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCTerminateNotify(baseRequest, baseResponse)
	case protobuffer_def.CMD_OFFER_CALL_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCOfferCall(baseRequest, baseResponse)
	case protobuffer_def.CMD_OFFER_ANSWER_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCOfferAnswer(baseRequest, baseResponse)
	case protobuffer_def.CMD_ICE_CALL_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCIceCall(baseRequest, baseResponse)
	case protobuffer_def.CMD_ICE_ANSWER_SIG_TO_SIG:
		return impl.GetSigServiceImpl().RPCIceAnswer(baseRequest, baseResponse)
	case protobuffer_def.CMD_CMD_MEDIA_TO_AUDIO:
		return impl.GetSigServiceImpl().RPCMediaToAudio(baseRequest, baseResponse)
	case protobuffer_def.CMD_CMD_KICK_USER:
		return impl.GetSigServiceImpl().RPCKickUser(baseRequest, baseResponse)
	default:
		baseResponse.Desc = "unkown cmd"
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_CMD //示知的指令
	}
	return nil
}
