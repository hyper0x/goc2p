package ctcp

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	DELIM byte = '\t'
)

var once sync.Once
var benchmarkServerAddr string = "127.0.0.1:8081"
var benchmarkListener TcpListener

func TestPrimeFuncs(t *testing.T) {
	t.Parallel()
	showLog := false
	serverAddr := "127.0.0.1:8080"
	t.Logf("Test tcp listener & sender (serverAddr=%s)... %s\n",
		serverAddr, generateRuntimeInfo())
	listener := generateTcpListener(serverAddr, showLog)
	if listener == nil {
		t.Fatalf("Listener startup failing! (addr=%s)!\n", serverAddr)
	}
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()
	if testing.Short() {
		multiSend(serverAddr, "SenderT", 1, (2 * time.Second), showLog)
	} else {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			multiSend(serverAddr, "SenderT1", 2, (2 * time.Second), showLog)
		}()
		go func() {
			defer wg.Done()
			multiSend(serverAddr, "SenderT2", 1, (2 * time.Second), showLog)
		}()
		wg.Wait()
	}
}

func BenchmarkPrimeFuncs(t *testing.B) {
	showLog := false
	t.Logf("Benchmark tcp listener & sender (serverAddr=%s)... %s\n",
		benchmarkServerAddr, generateRuntimeInfo())
	once.Do(startupListenerOnce)
	if benchmarkListener == nil {
		t.Errorf("Listener startup failing! (addr=%s)!\n", benchmarkServerAddr)
	}
	//if t.N == 1 {
	//	fmt.Printf("\nIterations (N): %d\n", t.N)
	//} else {
	//	fmt.Printf("Iterations (N): %d\n", t.N)
	//}
	for i := 0; i < t.N; i++ {
		multiSend(benchmarkServerAddr, "SenderB", 3, (2 * time.Second), showLog)
	}
	if benchmarkListener != nil {
		benchmarkListener.Close()
	}
}

func startupListenerOnce() {
	benchmarkListener = generateTcpListener(benchmarkServerAddr, false)
}

func generateTcpListener(serverAddr string, showLog bool) TcpListener {
	var listener TcpListener = NewTcpListener()
	var hasError bool
	if showLog {
		fmt.Printf("Start Listening at address %s ...\n", serverAddr)
	}
	err := listener.Init(serverAddr)
	if err != nil {
		hasError = true
		fmt.Errorf("Listener Init error: %s", err)
	}
	err = listener.Listen(requestHandler(showLog))
	if err != nil {
		hasError = true
		fmt.Errorf("Listener Listen error: %s", err)
	}
	if !hasError {
		return listener
	} else {
		if listener != nil {
			listener.Close()
		}
		return nil
	}
}

func multiSend(
	remoteAddr string,
	clientName string,
	number int,
	timeout time.Duration,
	showLog bool) {
	sender := NewTcpSender()
	if showLog {
		fmt.Printf("Initializing sender (%s) (remote address: %s, timeout: %d) ...", clientName, remoteAddr, timeout)
	}
	err := sender.Init(remoteAddr, timeout)
	if err != nil {
		fmt.Errorf("%s: Init Error: %s\n", clientName, err)
		return
	}
	if number <= 0 {
		number = 5
	}
	for i := 0; i < number; i++ {
		content := generateTestContent(fmt.Sprintf("%s-%d", clientName, i))
		if showLog {
			fmt.Printf("%s: Send content: '%s'\n", clientName, content)
		}
		err := sender.Send(content)
		if err != nil {
			fmt.Errorf("%s: Send Error: %s\n", clientName, err)
		}
		respChan := sender.Receive(DELIM)
		var resp TcpMessage
		timeoutChan := time.After(1 * time.Second)
		select {
		case resp = <-respChan:
		case <-timeoutChan:
			break
		}
		if err = resp.Err(); err != nil {
			fmt.Errorf("Sender: Receive Error: %s\n", err)
		} else {
			if showLog {
				respContent := resp.Content()
				fmt.Printf("%s: Received response: '%s'\n", clientName, respContent)
			}
		}
	}
	content := generateTestContent(fmt.Sprintf("%s-quit", clientName))
	if showLog {
		fmt.Printf("%s: Send content: '%s'\n", clientName, content)
	}
	err = sender.Send(content)
	if err != nil {
		fmt.Errorf("%s: Send Error: %s\n", clientName, err)
	}
	sender.Close()
}

func generateTestContent(content string) string {
	var respBuffer bytes.Buffer
	respBuffer.WriteString(strings.TrimSpace(content))
	respBuffer.WriteByte(DELIM)
	return respBuffer.String()
}

func requestHandler(showLog bool) func(conn net.Conn) {
	return func(conn net.Conn) {
		for {
			content, err := Read(conn, DELIM)
			if err != nil {
				fmt.Errorf("Listener Read error: %s", err)
			} else {
				if showLog {
					fmt.Printf("Listener: Received content: '%s'\n", content)
				}
				content = strings.TrimSpace(content)
				if strings.HasSuffix(content, "quit") {
					if showLog {
						fmt.Println("Listener: Quit!")
					}
					break
				}
				resp := generateTestContent(fmt.Sprintf("Resp: %s", content))
				n, err := Write(conn, resp)
				if err != nil {
					fmt.Errorf("Listener Write error: %s", err)
				}
				if showLog {
					fmt.Println("Listener: Send response: '%s' (n=%d)\n", resp, n)
				}
			}
		}
	}
}

func generateRuntimeInfo() string {
	return fmt.Sprintf("[GOMAXPROCS=%d, NUM_CPU=%d, NUM_GOROUTINE=%d]",
		runtime.GOMAXPROCS(-1), runtime.NumCPU(), runtime.NumGoroutine())
}
