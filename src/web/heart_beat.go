package web

import (
	"community-sig/config"
	"community-sig/logging"
	"community-sig/model"
	"time"
)

func CloseHeartBeat() {
	close(closeHeart)
}

func StartHeartBeat() {
	logging.Logger.Info("heart beat goroutine run")
	t1 := time.NewTimer(time.Duration(config.GetConf().Base.HeartCheck) * time.Second)
	defer func() {
		if err := recover(); err != nil {
			logging.Logger.Info("heart beat goroutine recover from panic ")
		}
		logging.Logger.Info("heart beat goroutine exit")
		t1.Stop()
		close(heartExit)
	}()
	for {
		select {
		case <-closeHeart:
			logging.Logger.Info("receive close heart msg, goroutine exited")
			return

		case deathClient := <-model.GetFromReconKick():
			logging.Logger.Info("get client from recon kick chan")
			model.GetClientMgr().DelClient(deathClient)
			continue
		case <-t1.C:
			deathIds := model.GetClientMgr().CheckHeart()
			hw := GetMsgHandler(model.WS_Offline_SYS)
			for _, id := range deathIds {
				var param interface{}
				hw.HandleMsg(id, param)
			}
			t1.Reset(time.Duration(config.GetConf().Base.HeartCheck) * time.Second)
			continue
		}

	}
}

func WaitHeartExit() <-chan struct{} {
	return heartExit
}
