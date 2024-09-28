package internal

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/tidwall/resp"
	"log/slog"
	"net"
	"os"
)

type Server struct {
	Config
	peers        map[*Peer]bool
	ln           net.Listener
	addPeerCh    chan *Peer
	removePeerCh chan *Peer
	messageCh    chan Message
	exitCh       chan struct{}

	Data cmap.ConcurrentMap[string, string]
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = ":" + os.Getenv("HOST_PORT")
	}

	return &Server{
		Config:       cfg,
		peers:        make(map[*Peer]bool),
		addPeerCh:    make(chan *Peer),
		removePeerCh: make(chan *Peer),
		messageCh:    make(chan Message),
		exitCh:       make(chan struct{}),
	}
}

func (s *Server) Start() error {
	// starting listening on tcp request
	ln, err := net.Listen("tcp", s.Config.ListenAddress)
	if err != nil {
		return err
	}

	// starting go routine loop to handle channels
	go s.loop()
	s.ln = ln
	slog.Info("GoCacher server running", "listenAddr", s.Config.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		// accepting connection
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "err", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.messageCh, s.removePeerCh)
	s.addPeerCh <- peer
	// handling message using peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func (s *Server) handleMessage(msg Message) error {

	//handle specific types of msg by command
	switch v := msg.cmd.(type) {
	case CommandHealth:
		if err := resp.NewWriter(msg.peer.conn).WriteString("OK"); err != nil {
			return err
		}
	case CommandSet:
		msg.peer.Data.Set(v.key, v.value)

		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil {
			return err
		}
	case CommandGet:
		val, ok := msg.peer.Data.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}

		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString(string(val)); err != nil {
			return err
		}
	case CommandRemove:
		msg.peer.Data.Remove(v.key)

		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil {
			return err
		}
	}

	return nil
}

// handling actions from channels
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.messageCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message eror", "err", err)
			}
		case <-s.exitCh:
			return
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.removePeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		}
	}
}
