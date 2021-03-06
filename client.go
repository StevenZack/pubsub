package pubsub

import (
	"errors"
	"fmt"

	"github.com/StevenZack/tools/strToolkit"
)

type Client struct {
	channelID string
	server    *Server
	id        int64
	receiver  chan interface{}
	log       bool
}

var autoIncrementId = &strToolkit.AutoIncrementID{}

func NewClient(s *Server) *Client {
	return &Client{
		server:   s,
		id:       autoIncrementId.Generate(),
		receiver: make(chan interface{}, 24),
	}
}

func (c *Client) SetLog(b bool) {
	c.log = b
}

func (c *Client) Sub(chanID string, listener func(interface{})) error {
	c.channelID = chanID
	if c.log {
		fmt.Println("client : sending entering")
	}
	c.server.entering <- c
	if c.log {
		fmt.Println("client : entering sent")
	}
	for msg := range c.receiver {
		listener(msg)
	}
	return nil
}

func (c *Client) UnSub() error {
	c.closeChannel()
	select {
	case c.server.leaving <- c:
		return nil
	default:
		return errors.New("server not running")
	}
}

func (c *Client) closeChannel() {
	select {
	case <-c.receiver:
		break
	default:
		close(c.receiver)
	}
}

func (c *Client) ID() int64 {
	return c.id
}
