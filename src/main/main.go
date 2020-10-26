package main

import (
	"community-sig/config"
	"community-sig/grpc_client"
	"community-sig/grpc_server"
	"community-sig/logging"
	"community-sig/web"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	//加载初始化.toml默认配置
	if err := config.LoadConfigAndSetDefault(); err != nil {
		panic(err.Error())
	}

	//初始化日志配置
	if err := logging.InitZap(&config.GetConf().LogConf); err != nil {
		panic("InitLogger:" + err.Error())
	}
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//启动web服务
	go web.Run()
	//启动心跳检测服务
	go web.StartHeartBeat()

	go grpc_server.StartComSigService(config.GetConf())

	grpc_client.StatusClient()
	select {
	case sign := <-c:
		logging.Logger.Info("catch stop signal, ", sign)
		fmt.Println("catch stop signal", sign)
		//服务退出
		web.Shutdown()
		web.CloseHeartBeat()
		grpc_server.StopComSigService()
		grpc_client.CloseStatusServiceClient()
		grpc_client.CloseSigClients()
		time.Sleep(time.Second * 3) //等待三秒
		os.Exit(0)
	case <-web.WaitHeartExit():
		fmt.Println("catch heart beat exit")
		//服务退出
		web.Shutdown()
		grpc_server.StopComSigService()
		grpc_client.CloseStatusServiceClient()
		grpc_client.CloseSigClients()
		time.Sleep(time.Second * 3) //等待三秒
		os.Exit(0)
	}

}
