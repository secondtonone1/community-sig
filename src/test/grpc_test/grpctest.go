package main

import (
	"community-sig/grpc_client"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	grpc_client.Start()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		fmt.Println("catch exit signal ")
	}

}
