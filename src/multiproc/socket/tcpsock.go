package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t'
)

var wg sync.WaitGroup

var logSn = 1

func printLog(format string, args ...interface{}) {
	fmt.Printf("%d: %s", logSn, fmt.Sprintf(format, args...))
	logSn++
}

func convertToInt32(str string) (int32, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		printLog(fmt.Sprintf("Parse Error: %s\n", err))
		return 0, errors.New(fmt.Sprintf("'%s' is not integer!", str))
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		printLog(fmt.Sprintf("Convert Error: The integer %s is too large/small.\n", num))
		return 0, errors.New(fmt.Sprintf("'%s' is not 32-bit integer!", num))
	}
	return int32(num), nil
}

func cbrt(param int32) float64 {
	return math.Cbrt(float64(param))
}

// 千万不要使用这个版本的read函数！
//func read(conn net.Conn) (string, error) {
//	reader := bufio.NewReader(conn)
//	readBytes, err := reader.ReadBytes(DELIMITER)
//	if err != nil {
//		return "", err
//	}
//	return string(readBytes[:len(readBytes)-1]), nil
//}

func read(conn net.Conn) (string, error) {
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

func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}

func serverGo() {
	var listener net.Listener
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		printLog("Listen Error: %s\n", err)
		return
	}
	defer listener.Close()
	printLog("Got listener for the server. (local address: %s)\n",
		listener.Addr())
	for {
		conn, err := listener.Accept() // 阻塞直至新连接到来
		if err != nil {
			printLog("Accept Error: %s\n", err)
		}
		printLog("Established a connection with a client application. (remote address: %s)\n",
			conn.RemoteAddr())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		wg.Done()
	}()
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		strReq, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printLog("The connection is closed by another side. (Server)\n")
			} else {
				printLog("Read Error: %s (Server)\n", err)
			}
			break
		}
		printLog("Received request: %s (Server)\n", strReq)
		i32Req, err := convertToInt32(strReq)
		if err != nil {
			n, err := write(conn, err.Error())
			if err != nil {
				printLog("Write Error (written %d bytes): %s (Server)\n", err)
			}
			printLog("Sent response (written %d bytes): %s (Server)\n", n, err)
			continue
		}
		f64Resp := cbrt(i32Req)
		respMsg := fmt.Sprintf("The cube root of %d is %f.", i32Req, f64Resp)
		n, err := write(conn, respMsg)
		if err != nil {
			printLog("Write Error: %s (Server)\n", err)
		}
		printLog("Sent response (written %d bytes): %s (Server)\n", n, respMsg)
	}
}

func clientGo(id int) {
	defer wg.Done()
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		printLog("Dial Error: %s (Client[%d])\n", err, id)
		return
	}
	defer conn.Close()
	printLog("Connected to server. (remote address: %s, local address: %s) (Client[%d])\n",
		conn.RemoteAddr(), conn.LocalAddr(), id)
	time.Sleep(200 * time.Millisecond)
	requestNumber := 5
	conn.SetDeadline(time.Now().Add(5 * time.Millisecond))
	for i := 0; i < requestNumber; i++ {
		i32Req := rand.Int31()
		n, err := write(conn, fmt.Sprintf("%d", i32Req))
		if err != nil {
			printLog("Write Error: %s (Client[%d])\n", err, id)
			continue
		}
		printLog("Sent request (written %d bytes): %d (Client[%d])\n", n, i32Req, id)
	}
	for j := 0; j < requestNumber; j++ {
		strResp, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printLog("The connection is closed by another side. (Client[%d])\n", id)
			} else {
				printLog("Read Error: %s (Client[%d])\n", err, id)
			}
			break
		}
		printLog("Received response: %s (Client[%d])\n", strResp, id)
	}
}

func main() {
	wg.Add(2)
	go serverGo()
	time.Sleep(500 * time.Millisecond)
	go clientGo(1)
	wg.Wait()
}
