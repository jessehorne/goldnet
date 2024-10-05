package main

import (
	"github.com/jessehorne/goldnet/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := "127.0.0.1:5555"
	s, err := server.NewServer(addr)
	if err != nil {
		panic(err)
	}
	s.Start()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	s.Stop()
}
