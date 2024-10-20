package server

import (
	"bufio"
	"fmt"
	"github.com/jessehorne/goldnet/internal/config"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/server/handlers"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
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
	Conf       *config.ServerConfig
}

func NewServer(conf *config.ServerConfig) (*Server, error) {
	listener, err := net.Listen("tcp", conf.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", conf.ServerAddress, err)
	}

	gs := game.NewGameState()
	go gs.RunGameLoop()

	return &Server{
		Conf:       conf,
		Listener:   listener,
		Shutdown:   make(chan struct{}),
		Connection: make(chan net.Conn),
		Logger:     log.New(os.Stdout, "[GoldNet] (Server) ", log.Ldate|log.Ltime),
		GameState:  gs,
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
	playerID := s.GameState.NextPlayerID()
	for {
		lenBytes := make([]byte, 8)
		_, err := io.ReadFull(reader, lenBytes)
		if err != nil {
			handlers.ServerUserDisconnectedHandler(s.GameState, playerID, conn, nil)
			break
		}

		msgLen := util.BytesToInt64(lenBytes)
		msgBytes := make([]byte, msgLen)
		_, err = io.ReadFull(reader, msgBytes)
		if err != nil {
			handlers.ServerUserDisconnectedHandler(s.GameState, playerID, conn, nil)
			break
		}

		msg := packets.Raw{}
		err = proto.Unmarshal(msgBytes, &msg)
		if err != nil {
			handlers.ServerUserDisconnectedHandler(s.GameState, playerID, conn, nil)
			break
		}

		// handle the packet and start over
		handler.Handle(playerID, conn, &msg, msgBytes)
	}
}

func (s *Server) Start() {
	s.Logger.Println("server started on", s.Conf.ServerAddress)
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
