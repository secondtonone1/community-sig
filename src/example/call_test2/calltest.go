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

//var addr = flag.String("addr", "localhost:9599", "http service address")

//var addr = flag.String("addr", "81.68.86.146:9499", "http service address")

var addr = flag.String("addr", "180.76.163.81:9499", "http service address")

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
	loginReq.UserId = "1025"
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

			req := &model.RequestStruct{}
			err = json.Unmarshal(message, req)
			if err != nil {
				log.Println("json unmarshal failed , error is ", err.Error())
				continue
			}
			fmt.Println("...................................")
			singleNotify := &model.SCCallSingleNotify{}
			if err := mapstructure.Decode(req.Data, singleNotify); err != nil {
				log.Println(" map to struct failed, err is ", err.Error())
				continue
			}

			fmt.Println("*********************************")

			log.Println("singleNotify.CallerId is", singleNotify.CallerId)
			log.Println("singleNotify.ChatRoomId is", singleNotify.ChatRoomId)
			log.Println("singleNotify.MediaType is", singleNotify.MediaType)

			//发送answer回复，同意还是拒绝
			answer := &model.CSAnswerSingle{}
			answer.AnswerId = "1025"
			answer.ChatRoomId = singleNotify.ChatRoomId
			answerReq := model.RequestStruct{}
			answerReq.Event = model.WS_CALL_SINGLE_ANSWER_CS
			answerReq.Data = answer
			senddata, err = json.Marshal(answerReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}

			//发送answer回复，同意还是拒绝
			/*
				refuse := &model.CSRefuseSingle{}
				refuse.AnswerId = "1025"
				refuse.ChatRoomId = singleNotify.ChatRoomId
				refuseReq := model.RequestStruct{}
				refuseReq.Event = model.WS_CALL_SINGLE_REFUSE_CS
				refuseReq.Data = refuse
				senddata, err = json.Marshal(refuseReq)
				err = c.WriteMessage(websocket.TextMessage, senddata)
				if err != nil {
					log.Println("write:", err)
					return
				}
			*/
			/*
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

				hangNotify := &model.SCHangupSingleNotify{}
				if err := mapstructure.Decode(req.Data, hangNotify); err != nil {
					log.Println(" map to struct failed, err is ", err.Error())
					continue
				}

				log.Println("hangNotify.ChatRoomId ", hangNotify.ChatRoomId)
			*/
			/*
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

				terminateNotify := &model.SCTerminateSingleNotify{}
				if err := mapstructure.Decode(req.Data, terminateNotify); err != nil {
					log.Println(" map to struct failed, err is ", err.Error())
					continue
				}

				log.Println("terminateNotify.ChatRoomId ", terminateNotify.ChatRoomId)
				log.Println("terminateNotify.CancelId ", terminateNotify.CancelId)

			*/

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

			offerCallNotify := &model.SCOfferCallNotify{}
			if err := mapstructure.Decode(req.Data, offerCallNotify); err != nil {
				log.Println(" map to struct failed, err is ", err.Error())
				continue
			}

			log.Println("offerCallNotify.ChatRoomId ", offerCallNotify.ChatRoomId)
			log.Println("offerCallNotify.CallerId ", offerCallNotify.CallerId)
			log.Println("offerCallNotify.Sdp ", offerCallNotify.Sdp)

			//发送offer answer
			offerAnswer := &model.CSOfferAnswer{}
			offerAnswer.AnswerId = "1025"
			offerAnswer.ChatRoomId = singleNotify.ChatRoomId
			offerAnswer.Sdp = "good answer sdp"
			offerAnswerReq := model.RequestStruct{}
			offerAnswerReq.Event = model.WS_OFFER_ANSWER_CS
			offerAnswerReq.Data = offerAnswer
			senddata, err = json.Marshal(offerAnswerReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}

			//ice 接收和answer ice
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

			iceCallNotify := &model.SCIceCallNotify{}
			if err := mapstructure.Decode(req.Data, iceCallNotify); err != nil {
				log.Println(" map to struct failed, err is ", err.Error())
				continue
			}

			log.Println("iceCallNotify.ChatRoomId ", iceCallNotify.ChatRoomId)
			log.Println("iceCallNotify.CallerId ", iceCallNotify.CallerId)
			log.Println("iceCallNotify.ice ", iceCallNotify.IceCandidate)

			//发送offer answer
			iceAnswer := &model.CSIceAnswer{}
			iceAnswer.AnswerId = "1025"
			iceAnswer.ChatRoomId = singleNotify.ChatRoomId
			iceCandidate, _ := json.Marshal(&model.CSTerminateSingle{ChatRoomId: iceCallNotify.ChatRoomId})
			iceAnswer.IceCandidate = string(iceCandidate)
			iceAnswerReq := model.RequestStruct{}
			iceAnswerReq.Event = model.WS_ICE_ANSWER_CS
			iceAnswerReq.Data = iceAnswer
			senddata, err = json.Marshal(iceAnswerReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}

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
			heartBeat.UserId = "1025"
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
			case <-time.After(time.Second * 10):
			}
			return
		}
	}
}
