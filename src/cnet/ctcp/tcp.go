package ctcp

import (
	"bufio"
	"bytes"
	"errors"
	"logging"
	"net"
	"sync"
	"time"
)

const (
	DELIMITER = '\t'
)

var logger logging.Logger = logging.NewSimpleLogger()

func Read(conn net.Conn, delim byte) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

func Write(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWriter(conn)
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

type AsyncTcpListener struct {
	listener net.Listener
	active   bool
	lock     *sync.Mutex
}

func (self *AsyncTcpListener) Init(addr string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		return nil
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	self.listener = ln
	self.active = true
	return nil
}

func (self *AsyncTcpListener) Listen(handler func(conn net.Conn)) error {
	if !self.active {
		return errors.New("Listen Error: Uninitialized listener!")
	}
	go func(active *bool) {
		for {
			if *active {
				return
			}
			conn, err := self.listener.Accept()
			if err != nil {
				logger.Errorf("Listener: Accept Request Error: %s\n", err)
				continue
			}
			go handler(conn)
		}
	}(&self.active)
	return nil
}

func (self *AsyncTcpListener) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.listener.Close()
		self.active = false
		return true
	} else {
		return false
	}
}

func (self *AsyncTcpListener) Addr() net.Addr {
	if self.active {
		return self.listener.Addr()
	} else {
		return nil
	}
}

func NewTcpListener() TcpListener {
	return &AsyncTcpListener{lock: new(sync.Mutex)}
}

type AsyncTcpSender struct {
	active bool
	lock   *sync.Mutex
	conn   net.Conn
}

func (self *AsyncTcpSender) Init(remoteAddr string, timeout time.Duration) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.active {
		conn, err := net.DialTimeout("tcp", remoteAddr, timeout)
		if err != nil {
			return err
		}
		self.conn = conn
		self.active = true
	}
	return nil
}

func (self *AsyncTcpSender) Send(content string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.active {
		return errors.New("Send Error: Uninitialized sender!")
	}
	_, err := Write(self.conn, content)
	return err
}

func (self *AsyncTcpSender) Receive(delim byte) <-chan TcpMessage {
	respChan := make(chan TcpMessage, 1)
	go func(conn net.Conn, ch chan<- TcpMessage) {
		content, err := Read(conn, delim)
		ch <- NewTcpMessage(content, err)
	}(self.conn, respChan)
	return respChan
}

func (self *AsyncTcpSender) Addr() net.Addr {
	if self.active {
		return self.conn.LocalAddr()
	} else {
		return nil
	}
}

func (self *AsyncTcpSender) RemoteAddr() net.Addr {
	if self.active {
		return self.conn.RemoteAddr()
	} else {
		return nil
	}
}

func (self *AsyncTcpSender) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.conn.Close()
		self.active = false
		return true
	} else {
		return false
	}
}

func NewTcpSender() TcpSender {
	return &AsyncTcpSender{lock: new(sync.Mutex)}
}
