package model

const (
	WS_Offline_SYS               = "ws_offline_sys"
	WS_Login_CS                  = "ws_login_cs"
	WS_LOGIN_SC                  = "ws_login_sc"
	WS_HEART_BEAT_CS             = "ws_heart_beat_cs"
	WS_CALL_SINGLE_CS            = "ws_call_single_cs"
	WS_CALL_SINGLE_SC            = "ws_call_single_sc"
	WS_CALL_SINGLE_NOTIFY        = "ws_call_single_notify_sc"
	WS_CALL_SINGLE_ANSWER_CS     = "ws_call_single_answer_cs"
	WS_CALL_SINGLE_ANSWER_NOTIFY = "ws_call_single_answer_notify_sc"
	WS_CALL_SINGLE_REFUSE_CS     = "ws_call_single_refuse_cs"
	WS_CALL_SINGLE_REFUSE_NOTIFY = "ws_call_single_refuse_notify_sc"
	WS_SINGLE_TERMINATE_CS       = "ws_single_terminate_cs"
	WS_SINGLE_TERMINATE_NoTIFY   = "ws_single_terminate_notify_sc"
	WS_SINGLE_HANGUP_CS          = "ws_single_hangup_cs"
	WS_SINGLE_HANGUP_NOTIFY      = "ws_single_hangup_notify_sc"
	WS_OFFER_CALL_CS             = "ws_offer_call_cs"
	WS_OFFER_CALL_NOTIFY         = "ws_offer_call_notify_sc"
	WS_OFFER_ANSWER_CS           = "ws_offer_answer_cs"
	WS_OFFER_ANSWER_NOTIFY       = "ws_offer_answer_notify_sc"
	WS_ICE_CALL_CS               = "ws_ice_call_cs"
	WS_ICE_CALL_NOTIFY           = "ws_ice_call_notify_sc"
	WS_ICE_ANSWER_CS             = "ws_ice_answer_cs"
	WS_ICE_ANSWER_NOTIFY         = "ws_ice_answer_notify_sc"
	WS_MEDIA_TO_AUDIO_CS         = "ws_media_to_audio_cs"
	WS_MEDIA_TO_AUDIO_NOTIFY     = "ws_media_to_audio_notify_sc"
	WS_FORCE_TERMINATE_NOTIFY    = "ws_force_terminate_notify_sc"
	WS_CALL_MULT_CS              = "ws_call_multiple_cs"
	WS_CALL_MULT_SC              = "ws_call_multiple_sc"
	WS_CALL_MULT_NOTIFY          = "ws_call_multiple_notify_sc"
	WS_CALL_MULT_ANSWER_CS       = "ws_call_multiple_answer_cs"
	WS_CALL_MULT_ANSWER_SC       = "ws_call_multiple_answer_sc"
	WS_CALL_MULT_ANSWER_NOTIFY   = "ws_call_multiple_answer_notify_sc"
	WS_MULT_OTHER_ACCEPT_NOTIFY  = "ws_multiple_other_accept_notify_sc"
	WS_CALL_MULT_REFUSE_CS       = "ws_call_multiple_refuse_cs"
	WS_CALL_MULT_REFUSE_NOTIFY   = "ws_call_multiple_refuse_notify_sc"
	WS_CALL_MULT_HANGUP_CS       = "ws_multiple_hangup_cs"
	WS_CALL_MULT_HANGUP_NOTIFY   = "ws_multiple_hangup_notify_sc"
)

type CSLogin struct {
	UserId   string   `json:"userId" mapstructure:"userId"`
	UserName string   `json:"userName" mapstructure:"userName"`
	Phone    string   `json:"phone" mapstructure:"phone"`
	Avator   string   `json:"userAvatar" mapstructure:"userAvatar"`
	RoomList []string `json:"roomList" mapstructure:"roomList"`
}

type SCLogin struct {
	UserId string `json:"userId"`
}

type CSOffline struct {
}

type CSHeartBeat struct {
	UserId string `json:"userId" mapstructure:"userId"`
}

type CSCallSingle struct {
	CallerId    string `json:"callerId" mapstructure:"callerId"`
	AnswerId    string `json:"answerId" mapstructure:"answerId"`
	MediaType   string `json:"mediaType" mapstructure:"mediaType"`
	DeviceModel string `json:"deviceModel" mapstructure:"deviceModel"`
}

type SCCallSingle struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCCallSingleNotify struct {
	CallerId     string `json:"callerId" mapstructure:"callerId"`
	CallerName   string `json:"callerName" mapstructure:"callerName"`
	CallerPhone  string `json:"callerPhone" mapstructure:"callerPhone"`
	CallerAvator string `json:"callerAvatar" mapstructure:"callerAvatar"`
	MediaType    string `json:"mediaType" mapstructure:"mediaType"`
	ChatRoomId   string `json:"chatRoomId" mapstructure:"chatRoomId"`
	DeviceModel  string `json:"deviceModel" mapstructure:"deviceModel"`
}

type CSAnswerSingle struct {
	AnswerId    string `json:"answerId" mapstructure:"answerId"`
	ChatRoomId  string `json:"chatRoomId" mapstructure:"chatRoomId"`
	DeviceModel string `json:"deviceModel" mapstructure:"deviceModel"`
}

type SCAnswerSingleNotify struct {
	AnswerId   string `json:"answerId" mapstructure:"answerId"`
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`

	AnswerName   string `json:"answerName" mapstructure:"answerName"`
	AnswerPhone  string `json:"answerPhone" mapstructure:"answerPhone"`
	AnswerAvator string `json:"answerAvatar" mapstructure:"answerAvatar"`
	DeviceModel  string `json:"deviceModel" mapstructure:"deviceModel"`
}

type CSRefuseSingle struct {
	AnswerId   string `json:"answerId" mapstructure:"answerId"`
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCRefuseSingleNotify struct {
	AnswerId   string `json:"answerId" mapstructure:"answerId"`
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type CSTerminateSingle struct {
	CancelId   string `json:"cancelId" mapstructure:"cancelId"`
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCTerminateSingleNotify struct {
	CancelId   string `json:"cancelId" mapstructure:"cancelId"`
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type CSHangupSingle struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCHangupSingleNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type CSOfferCall struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
	CallerId   string `json:"callerId" mapstructure:"callerId"`
	Sdp        string `json:"sdp" mapstructure:"sdp"`
}

type SCOfferCallNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
	CallerId   string `json:"callerId" mapstructure:"callerId"`
	Sdp        string `json:"sdp" mapstructure:"sdp"`
}

type CSOfferAnswer struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
	AnswerId   string `json:"answerId" mapstructure:"answerId"`
	Sdp        string `json:"sdp" mapstructure:"sdp"`
}

type SCOfferAnswerNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
	AnswerId   string `json:"answerId" mapstructure:"answerId"`
	Sdp        string `json:"sdp" mapstructure:"sdp"`
}

type CSIceCall struct {
	ChatRoomId   string `json:"chatRoomId" mapstructure:"chatRoomId"`
	CallerId     string `json:"callerId" mapstructure:"callerId"`
	IceCandidate string `json:"iceCandidate" mapstructure:"iceCandidate"`
}

type SCIceCallNotify struct {
	ChatRoomId   string `json:"chatRoomId" mapstructure:"chatRoomId"`
	CallerId     string `json:"callerId" mapstructure:"callerId"`
	IceCandidate string `json:"iceCandidate" mapstructure:"iceCandidate"`
}

type CSIceAnswer struct {
	ChatRoomId   string `json:"chatRoomId" mapstructure:"chatRoomId"`
	AnswerId     string `json:"answerId" mapstructure:"answerId"`
	IceCandidate string `json:"iceCandidate" mapstructure:"iceCandidate"`
}

type SCIceAnswerNotify struct {
	ChatRoomId   string `json:"chatRoomId" mapstructure:"chatRoomId"`
	AnswerId     string `json:"answerId" mapstructure:"answerId"`
	IceCandidate string `json:"iceCandidate" mapstructure:"iceCandidate"`
}

type CSMediaToAudio struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
	ConverId   string `json:"converId" mapstructure:"converId"`
}

type SCMediaToAudioNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
	ConverId   string `json:"converId" mapstructure:"converId"`
}

type SCForceTerminateNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type CSCallMul struct {
	CallerId    string `json:"callerId" mapstructure:"callerId"`
	RoomId      string `json:"roomId" mapstructure:"roomId"`
	MediaType   string `json:"mediaType" mapstructure:"mediaType"`
	DeviceModel string `json:"deviceModel" mapstructure:"deviceModel"`
}

type SCCallMul struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCCallMulNotify struct {
	CallerId    string `json:"callerId" mapstructure:"callerId"`
	ChatRoomId  string `json:"chatRoomId" mapstructure:"chatRoomId"`
	MediaType   string `json:"mediaType" mapstructure:"mediaType"`
	DeviceModel string `json:"deviceModel" mapstructure:"deviceModel"`
}

type CSCallMulAnswer struct {
	AnswerId    string `json:"answerId" mapstructure:"answerId"`
	ChatRoomId  string `json:"chatRoomId" mapstructure:"chatRoomId"`
	DeviceModel string `json:"deviceModel" mapstructure:"deviceModel"`
}

type SCCallMulAnswer struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCCallMulAnswerNotify struct {
	ChatRoomId  string `json:"chatRoomId" mapstructure:"chatRoomId"`
	DeviceModel string `json:"deviceModel" mapstructure:"deviceModel"`
}

type SCMulOtherAcceptNotify struct {
}

type CSCallMulRefuse struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCCallMulRefuseNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type CSMulHangup struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}

type SCMulHangupNotify struct {
	ChatRoomId string `json:"chatRoomId" mapstructure:"chatRoomId"`
}
