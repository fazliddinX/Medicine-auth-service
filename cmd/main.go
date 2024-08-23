package main

import (
	"auth-service/cmd/servers"
	"auth-service/pkg/config"
	"auth-service/pkg/logs"
	"log"
)

func main() {
	log.Println("Starting server")
	logger := logs.InitLogger()

	cfg := config.Load()

	server := servers.NewServer(logger, cfg)

	wait := make(chan int)

	go server.RunGinServer()
	go server.RunGRPCServer()

	<-wait
}
