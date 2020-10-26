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

var addr = flag.String("addr", "localhost:9599", "http service address")

//var addr = flag.String("addr", "81.68.86.146:9599", "http service address")

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

			//发送挂断请求
			answer := &model.CSHangupSingle{}
			answer.ChatRoomId = anserNotify.ChatRoomId
			answerReq := model.RequestStruct{}
			answerReq.Event = model.WS_SINGLE_HANGUP_CS
			answerReq.Data = answer
			senddata, err = json.Marshal(answerReq)
			err = c.WriteMessage(websocket.TextMessage, senddata)
			if err != nil {
				log.Println("write:", err)
				return
			}

			message = make([]byte, 1024)
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

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
