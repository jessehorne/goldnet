package server

import (
	"bufio"
	"fmt"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/server/handlers"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type Server struct {
	Addr       string
	WG         sync.WaitGroup
	Listener   net.Listener
	Shutdown   chan struct{}
	Connection chan net.Conn
	Logger     *log.Logger
	GameState  *game.GameState
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", address, err)
	}

	return &Server{
		Addr:       address,
		Listener:   listener,
		Shutdown:   make(chan struct{}),
		Connection: make(chan net.Conn),
		Logger:     log.New(os.Stdout, "[GoldNet] (Server) ", log.Ldate|log.Ltime),
		GameState:  game.NewGameState(),
	}, nil
}

func (s *Server) AcceptConnections() {
	defer s.WG.Done()

	for {
		select {
		case <-s.Shutdown:
			return
		default:
			conn, err := s.Listener.Accept()
			if err != nil {
				continue
			}
			s.Connection <- conn
		}
	}
}

func (s *Server) HandleConnections(handler *handlers.PacketHandler) {
	defer s.WG.Done()

	for {
		select {
		case <-s.Shutdown:
			return
		case conn := <-s.Connection:
			go s.HandleConnection(conn, handler)
		}
	}
}

func (s *Server) HandleConnection(conn net.Conn, handler *handlers.PacketHandler) {
	defer conn.Close()
	s.Logger.Println("connection made: ", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

	playerID := int64(len(s.GameState.Players))

	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				handlers.ServerUserDisconnectedHandler(s.GameState, playerID, conn, bytes)
				break
			} else {
				break
			}
		}
		handler.Handle(playerID, conn, bytes)
	}
}

func (s *Server) Start() {
	s.Logger.Println("server started on", s.Addr)
	s.WG.Add(2)

	handler := handlers.NewPacketHandler(s.GameState)

	go s.AcceptConnections()
	go s.HandleConnections(handler)
}

func (s *Server) Stop() {
	s.Logger.Println("shutting down")
	close(s.Shutdown)
	s.Listener.Close()

	done := make(chan struct{})
	go func() {
		s.WG.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.Logger.Println("shutdown complete")
		return
	}
}
