package model

import "community-sig/constants"

type ResponseCode struct {
	Code constants.ResponseCodeType `json:"code"`
	Desc string                     `json:"desc"`
}

type MessageStruct struct {
	Data  interface{} `json:"data"`
	Event string      `json:"event"`
}

type ResponseStruct struct {
	ResponseCode
	MessageStruct
}

type RequestStruct struct {
	MessageStruct
}
