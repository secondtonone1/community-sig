// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"community-sig/model"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/goinggo/mapstructure"
	"github.com/gorilla/websocket"
)

//var addr = flag.String("addr", "localhost:6699", "http service address")

//var addr = flag.String("addr", "81.68.86.146:9699", "http service address")
var addr = flag.String("addr", "180.76.163.81:9699", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "v1/wsmsg"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	loginReq := &model.CSLogin{}
	loginReq.Avator = "1Avatar"
	loginReq.Phone = "1phone"
	loginReq.UserId = "1024"
	loginReq.UserName = "102name"
	loginReq.RoomList = []string{"room102", "room101"}
	loginreq := model.RequestStruct{}
	loginreq.Event = model.WS_Login_CS
	loginreq.Data = loginReq
	senddata, err := json.Marshal(loginreq)
	err = c.WriteMessage(websocket.TextMessage, senddata)
	if err != nil {
		log.Println("write:", err)
		return
	}

	singleCall := &model.CSCallSingle{}
	singleCall.CallerId = "1024"
	singleCall.AnswerId = "1025"
	singleCall.MediaType = "audio"

	singleReq := model.RequestStruct{}
	singleReq.Event = model.WS_CALL_SINGLE_CS
	singleReq.Data = singleCall
	senddata, err = json.Marshal(singleReq)
	err = c.WriteMessage(websocket.TextMessage, senddata)
	if err != nil {
		log.Println("write:", err)
		return
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			message = make([]byte, 1024)
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			message = make([]byte, 1024)
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			req := &model.RequestStruct{}
			err = json.Unmarshal(message, req)
			if err != nil {
				log.Println("json unmarshal failed , error is ", err.Error())
				continue
			}
			fmt.Println("...................................")
			anserNotify := &model.SCAnswerSingleNotify{}
			if err := mapstructure.Decode(req.Data, anserNotify); err != nil {
				log.Println(" map to struct failed, err is ", err.Error())
				continue
			}

			log.Println("anserNotify.AnswerId ", anserNotify.AnswerId)
			log.Println("anserNotify.ChatRoomId ", anserNotify.ChatRoomId)

			// message = make([]byte, 1024)
			// _, message, err = c.ReadMessage()
			// if err != nil {
			// 	log.Println("read:", err)
			// 	return
			// }
			// log.Printf("recv: %s", message)

			// req = &model.RequestStruct{}
			// err = json.Unmarshal(message, req)
			// if err != nil {
			// 	log.Println("json unmarshal failed , error is ", err.Error())
			// 	continue
			// }
			// fmt.Println("...................................")
			// refuseNotify := &model.SCRefuseSingleNotify{}
			// if err := mapstructure.Decode(req.Data, refuseNotify); err != nil {
			// 	log.Println(" map to struct failed, err is ", err.Error())
			// 	continue
			// }

			// log.Println("refuseNotify.AnswerId ", refuseNotify.AnswerId)
			// log.Println("refuseNotify.ChatRoomId ", refuseNotify.ChatRoomId)

			/*
				hangup := &model.CSHangupSingle{}
				hangup.ChatRoomId = anserNotify.ChatRoomId
				hangupReq := model.RequestStruct{}
				hangupReq.Event = model.WS_SINGLE_HANGUP_CS
				hangupReq.Data = hangup
				senddata, err = json.Marshal(hangupReq)
				err = c.WriteMessage(websocket.TextMessage, senddata)
				if err != nil {
					log.Println("write:", err)
					return
				}
			*/
			/*
				terminate := &model.CSTerminateSingle{}
				terminate.ChatRoomId = anserNotify.ChatRoomId
				terminate.CancelId = "1024"
				terminateReq := model.RequestStruct{}
				terminateReq.Event = model.WS_SINGLE_TERMINATE_CS
				terminateReq.Data = terminate
				senddata, err = json.Marshal(terminateReq)
				err = c.WriteMessage(websocket.TextMessage, senddata)
				if err != nil {
					log.Println("write:", err)
					return
				}
			*/

			offercall := &model.CSOfferCall{}
			offercall.ChatRoomId = anserNotify.ChatRoomId
			offercall.CallerId = "1024"
			offercall.Sdp = "hello caller sdp"
			offercallReq := model.RequestStruct{}
			offercallReq.Event = model.WS_OFFER_CALL_CS
			offercallReq.Data = offercall
			senddata, err = json.Marshal(offercallReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}

			//接收offer
			message = make([]byte, 1024)
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			req = &model.RequestStruct{}
			err = json.Unmarshal(message, req)
			if err != nil {
				log.Println("json unmarshal failed , error is ", err.Error())
				continue
			}
			fmt.Println("...................................")
			offerAnswerNotify := &model.SCOfferAnswerNotify{}
			if err := mapstructure.Decode(req.Data, offerAnswerNotify); err != nil {
				log.Println(" map to struct failed, err is ", err.Error())
				continue
			}

			log.Println("offerAnswerNotify.AnswerId ", offerAnswerNotify.AnswerId)
			log.Println("offerAnswerNotify.ChatRoomId ", offerAnswerNotify.ChatRoomId)
			log.Println("offerAnswerNotify.Sdp ", offerAnswerNotify.Sdp)

			icecall := &model.CSIceCall{}
			icecall.ChatRoomId = anserNotify.ChatRoomId
			icecall.CallerId = "1024"
			iceCandidate, _ := json.Marshal(&model.CSTerminateSingle{ChatRoomId: anserNotify.ChatRoomId})
			icecall.IceCandidate = string(iceCandidate)

			icecallReq := model.RequestStruct{}
			icecallReq.Event = model.WS_ICE_CALL_CS
			icecallReq.Data = icecall
			senddata, err = json.Marshal(icecallReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}

			//接收ICE
			message = make([]byte, 1024)
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			req = &model.RequestStruct{}
			err = json.Unmarshal(message, req)
			if err != nil {
				log.Println("json unmarshal failed , error is ", err.Error())
				continue
			}
			fmt.Println("...................................")
			iceAnswerNotify := &model.SCIceAnswerNotify{}
			if err := mapstructure.Decode(req.Data, iceAnswerNotify); err != nil {
				log.Println(" map to struct failed, err is ", err.Error())
				continue
			}

			log.Println("iceAnswerNotify.AnswerId ", iceAnswerNotify.AnswerId)
			log.Println("iceAnswerNotify.ChatRoomId ", iceAnswerNotify.ChatRoomId)
			log.Println("iceAnswerNotify.Ice ", iceAnswerNotify.IceCandidate)

		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:

			heartBeat := &model.CSHeartBeat{}
			heartBeat.UserId = "1024"
			heartReq := model.RequestStruct{}
			heartReq.Event = model.WS_HEART_BEAT_CS
			heartReq.Data = heartBeat
			senddata, err = json.Marshal(heartReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}
			fmt.Println("test .....")
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
