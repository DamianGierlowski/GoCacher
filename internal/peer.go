package internal

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
)

type Peer struct {
	conn      net.Conn
	removeCh  chan *Peer
	messageCh chan Message
	Data      cmap.ConcurrentMap[string, string]
}

func NewPeer(conn net.Conn, messageCh chan Message, removeCh chan *Peer) *Peer {
	return &Peer{
		conn:      conn,
		messageCh: messageCh,
		removeCh:  removeCh,
		Data:      cmap.New[string](),
	}
}

func (p *Peer) readLoop() error {
	var rd = resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.removeCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var cmd Command

		if v.Type() == resp.Array {
			// fetching command from raw message
			rawCMD := v.Array()[0]
			switch rawCMD.String() {
			case HealthCommand:
				cmd = CommandHealth{}
			case SetCommand:
				cmd = CommandSet{
					key:   v.Array()[1].String(),
					value: v.Array()[2].String(),
				}
			case GetCommand:
				cmd = CommandGet{
					key: v.Array()[1].String(),
				}
			case RemoveCommand:
				cmd = CommandRemove{
					key: v.Array()[1].String(),
				}
			default:
				fmt.Println("got this unhandled command", rawCMD)
			}

		}

		//passing command to message channel
		p.messageCh <- Message{
			cmd:  cmd,
			peer: p,
		}
	}

	return nil
}
