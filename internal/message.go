package internal

type Message struct {
	cmd  Command
	peer *Peer
}
