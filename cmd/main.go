package main

import (
	"auth-service/cmd/servers"
	"auth-service/pkg/config"
	"auth-service/pkg/logs"
)

func main() {
	logger := logs.InitLogger()

	cfg := config.Load()

	server := servers.NewServer(logger, cfg)

	wait := make(chan int)

	go server.RunGinServer()
	go server.RunGRPCServer()

	<-wait
}
