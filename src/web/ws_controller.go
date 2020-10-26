package web

import (
	"bytes"
	"community-sig/logging"
	"community-sig/model"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   8192, //指定读缓存区大小
	WriteBufferSize:  8192, // 指定写缓存区大小
	HandshakeTimeout: 5 * time.Second,
	// 检测请求来源
	CheckOrigin: func(r *http.Request) bool {
		if r.Method != "GET" {
			fmt.Println("method is not GET")
			return false
		}
		return true
	},
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

/*
websocket 处理*/

func wsHandler(c *gin.Context) {

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logging.Logger.Info("websocket upgrade failed, err is ", err)
		return
	}
	clientInst := model.InitClient(ws)
	defer func() {
		//recover恢复
		if err := recover(); err != nil {
			logging.Logger.Info("web socket logic goroutine recover from panic, err is ", err)
		}

		alive := model.GetClientMgr().IsClientAlive(clientInst)
		if !alive {
			//被踢会走入这个逻辑
			return
		}
		//主动退出会走入这个逻辑
		logging.Logger.Info("web socket normal closed")
		hw := GetMsgHandler(model.WS_Offline_SYS)
		var param interface{}
		hw.HandleMsg(clientInst, param)
		model.GetClientMgr().DelClient(clientInst)
	}()

	model.GetClientMgr().AddClient(clientInst)

	for {

		//websocket接受信息
		_, message, err := ws.ReadMessage()
		if err != nil {
			logging.Logger.Info("message read failed, maybe peer closed")
			break
		}

		logging.Logger.Info("receive msg is ", string(message))
		req := &model.RequestStruct{}
		err = json.Unmarshal(message, req)
		if err != nil {
			logging.Logger.Info("json unmarshal failed , error is ", err.Error())
			continue
		}

		hw := GetMsgHandler(req.Event)
		if hw == nil {
			logging.Logger.Infof("msg %s isn't registered", req.Event)
			continue
		}

		//logging.Logger.Info("req.Data is ", req.Data)
		//logging.Logger.Info("req.Data type is ", reflect.TypeOf(req.Data).String())
		//处理逻辑消息
		err = hw.HandleMsg(clientInst, req.Data)

		if err != nil {
			logging.Logger.Infof("handle web msg[%s] failed, error is %s", req.Event, err.Error())
			break
		}

	}
}
