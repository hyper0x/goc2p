package ctcp

import (
	"net"
	"time"
)

type TcpMessage struct {
	content string
	err     error
}

func (self TcpMessage) Content() string {
	return self.content
}

func (self TcpMessage) Err() error {
	return self.err
}

func NewTcpMessage(content string, err error) TcpMessage {
	return TcpMessage{content: content, err: err}
}

type TcpListener interface {
	Init(addr string) error
	Listen(handler func(conn net.Conn)) error
	Close() bool
	Addr() net.Addr
}

type TcpSender interface {
	Init(remoteAddr string, timeout time.Duration) error
	Send(content string) error
	Receive(delim byte) <-chan TcpMessage
	Close() bool
	Addr() net.Addr
	RemoteAddr() net.Addr
}
