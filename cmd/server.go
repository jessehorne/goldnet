package main

import (
	"github.com/jessehorne/goldnet/internal/config"
	"github.com/jessehorne/goldnet/internal/server"
	"github.com/jessehorne/goldnet/internal/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := config.NewServerConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}

	util.PerlinInit(conf.WorldSeed)

	s, err := server.NewServer(conf)
	if err != nil {
		panic(err)
	}
	s.Start()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	s.Stop()
}
