package constants

import (
	"community-sig/config"
	"errors"
)

type ResponseCodeType int

const (
	ResponseCodeSuccess                 ResponseCodeType = 0
	ResponseCodeBodyIsNull              ResponseCodeType = 2
	ResponseCodeJsonParsError           ResponseCodeType = 3
	ResponseCodeLessParamError          ResponseCodeType = 4
	ResponseCodeServerError             ResponseCodeType = 5
	ResponseCodeParamError              ResponseCodeType = 6
	ResponseFail                        ResponseCodeType = 7
	ResponseHaveDataNoAllowDel          ResponseCodeType = 8
	ResponseCodeAuthError               ResponseCodeType = 403
	ResponseCodeUserError               ResponseCodeType = 9
	ResponseCodeOnlineError             ResponseCodeType = 10
	ResponseCodeRoomError               ResponseCodeType = 11
	ResponseUserBusyError               ResponseCodeType = 12
	ResponseOtherAnswer                 ResponseCodeType = 13 //其他用户接听了
	ResponseUserNotOnline               ResponseCodeType = 14
	ResponseRoomUserALLOff              ResponseCodeType = 15
	ResponseMulChatRoomError            ResponseCodeType = 16
	ResponseMulChatAnswerError          ResponseCodeType = 17
	ResponseCodeLoginFailed             ResponseCodeType = 18 //登录失败
	ResponseCodeRpcGetUserFailed        ResponseCodeType = 19
	ResopnseCodeRpcCreateChatRoomFailed ResponseCodeType = 20
	ResponseCodeRpcGetChatRoomFailed    ResponseCodeType = 21
)

func (p ResponseCodeType) String() string {
	switch p {
	case ResponseCodeSuccess:
		return "success"
	case ResponseCodeBodyIsNull:
		return "body is null"
	case ResponseCodeJsonParsError:
		return "json parse error"
	case ResponseCodeLessParamError:
		return "less param"
	case ResponseCodeServerError:
		return "server error"
	case ResponseCodeParamError:
		return "param error"
	case ResponseFail:
		return "fail"
	case ResponseHaveDataNoAllowDel:
		return "have data, not allow delete"
	case ResponseCodeAuthError:
		return "check auth fail"
	case ResponseCodeRoomError:
		return "room not exists"
	case ResponseUserBusyError:
		return "user is busy now"
	case ResponseUserNotOnline:
		return "user is not online"
	case ResponseRoomUserALLOff:
		return "all room user offline "
	case ResponseMulChatRoomError:
		return "mul chat room id invalid error "
	case ResponseCodeLoginFailed:
		return "rpc request login failed"
	case ResponseCodeRpcGetUserFailed:
		return "rpc get user data failed"
	case ResponseCodeRpcGetChatRoomFailed:
		return "rpc get chat room failed"
	default:
		return "unknown code desc"
	}
}

var (
	ErrMap2Struct           = errors.New("map to struct failed ")
	ErrUserNotFound         = errors.New("user not found")
	ErrCallSingle           = errors.New("user call single error")
	ErrJsonMarshal          = errors.New("json marshal failed")
	ErrSendData             = errors.New("send data failed")
	ErrClient               = errors.New("client has changed")
	ErrHangupIsNotCaller    = errors.New("hang up is not caller")
	ErrTerminateUserInvalid = errors.New("terminate user is invalid ")
	ErrOfferCallerInvalid   = errors.New("offer caller is invalid ")
	ErrOfferAnswerInvalid   = errors.New("offer answer is invalid ")
	ErrIceCallerInvalid     = errors.New("Ice caller invalid")
	ErrIceAnswerInvalid     = errors.New("Ice answer invalid")
	ErrMediaToAudioInvalid  = errors.New("media to audio invalid")
	ErrChatRoomInvalid      = errors.New("room id invalid ")
	ErrUserInvalid          = errors.New("user id invalid")
	ErrCallMul              = errors.New("user call mul error")
	ErrRpcCreateChatRoom    = errors.New("rpc create chat room failed")
	ErrRpcGetChatRoom       = errors.New("rpc get chat room failed ")
	ErrRpcAddrEmpty         = errors.New("rpc addr is empty")
)

const (
	User_Idle = 0
	User_Busy = 1
)

const (
	Recon_Chan_Size = 1024
)

func GetStatusRpcClientName() string {
	return "status-client-" + config.GetConf().Base.GRPCAddr
}
